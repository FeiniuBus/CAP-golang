package cap

type IPublish interface {
	Publish(keyName, content string) error
}

type QueueExecutorPublishBase struct {
	IQueueExecutor
	IPublish

	StateChanger IStateChanger
}

func (this *QueueExecutorPublishBase) Execute(connection IStorageConnection, feched IFetchedMessage) error {
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

	err = this.Publish(message.Name, message.Content)
	
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

func (this *QueueExecutorPublishBase) UpdateMessageForRetry(message *CapPublishedMessage, connection IStorageConnection) (bool, error) {
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