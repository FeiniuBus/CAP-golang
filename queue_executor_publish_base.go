package cap

type QueueExecutorPublish struct {
	IQueueExecutor
	StateChanger    IStateChanger
	PublishDelegate IPublishDelegate
}

func NewQueueExecutorPublish(stateChanger IStateChanger, delegate IPublishDelegate) *QueueExecutorPublish {
	return &QueueExecutorPublish{
		StateChanger:    stateChanger,
		PublishDelegate: delegate,
	}
}

func (this *QueueExecutorPublish) Execute(connection IStorageConnection, feched IFetchedMessage) error {
	message, err := connection.GetPublishedMessage(feched.GetMessageId())
	if err != nil {
		return err
	}

	transaction, err := connection.CreateTransaction()
	if err != nil {
		return err
	}

	err = this.StateChanger.ChangePublishedMessage(message, NewProcessingState(), transaction)
	if err != nil {
		return err
	}

	err = this.PublishDelegate.Publish(message.Name, message.Content)

	var newState IState = nil
	if err != nil {
		shouldRetry, err := this.UpdateMessageForRetry(message, connection)
		if err != nil {
			return err
		}

		if shouldRetry {
			newState = NewScheduledState()
		} else {
			newState = NewFailedState()
		}
	} else {
		newState = NewSucceededState()
	}

	err = this.StateChanger.ChangePublishedMessage(message, newState, transaction)
	if err != nil {
		return err
	}

	err = feched.RemoveFromQueue()
	if err != nil {
		return err
	}

	return nil
}

func (this *QueueExecutorPublish) UpdateMessageForRetry(message *CapPublishedMessage, connection IStorageConnection) (bool, error) {
	retryBehavior := DefaultRetry

	message.Retries = message.Retries + 1
	if message.Retries >= int(retryBehavior.RetryCount) {
		return false, nil
	}

	message.ExpiresAt = message.Added + int(retryBehavior.RetryIn(int32(message.Retries)))

	transaction, err := connection.CreateTransaction()
	if err != nil {
		return false, err
	}
	defer transaction.Dispose()

	err = transaction.UpdatePublishedMessage(message)
	if err != nil {
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
