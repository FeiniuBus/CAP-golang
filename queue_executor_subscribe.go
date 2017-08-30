package cap

type QueueExecutorSubscribe struct {
	IQueueExecutor
	Register *CallbackRegister
}

func NewQueueExecutorSubscribe(register *CallbackRegister) *QueueExecutorSubscribe {
	return &QueueExecutorSubscribe{
		Register: register,
	}
}

func (this *QueueExecutorSubscribe) Execute(connection IStorageConnection, feched IFetchedMessage) error {
	message, err := connection.GetReceivedMessage(feched.GetMessageId())
	if err != nil {
		return err
	}

	stateChanger := NewStateChanger()
	transaction, err := connection.CreateTransaction()
	if err != nil {
		return err
	}

	defer transaction.Dispose()

	err = stateChanger.ChangeReceivedMessageState(message, NewProcessingState(), transaction)
	if err != nil {
		return err
	}

	var newState IState
	err = this.executeSubscribeAsync(message)
	if err != nil {
		shouldRetry, err := this.updateMessageForRetryAsync(message, connection)
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

	transaction, err = connection.CreateTransaction()
	if err != nil {
		return err
	}

	err = stateChanger.ChangeReceivedMessageState(message, newState, transaction)
	if err != nil {
		return err
	}

	err = feched.RemoveFromQueue()
	if err != nil {
		return err
	}

	return nil
}

func (this *QueueExecutorSubscribe) executeSubscribeAsync(receivedMessage *CapReceivedMessage) error {
	cb, err := this.Register.Get(receivedMessage.Group, receivedMessage.Name)
	if err != nil {
		return nil
	}
	return cb.Handle(receivedMessage)
}

func (this *QueueExecutorSubscribe) updateMessageForRetryAsync(receivedMessage *CapReceivedMessage, conn IStorageConnection) (bool, error) {
	retryBehavior := DefaultRetry
	receivedMessage.Retries = receivedMessage.Retries + 1
	if receivedMessage.Retries >= int(retryBehavior.RetryCount) {
		return false, nil
	}

	due := receivedMessage.Added + int64(retryBehavior.RetryIn(int32(receivedMessage.Retries)))
	receivedMessage.ExpiresAt = int(due)

	transaction, err := conn.CreateTransaction()
	if err != nil {
		return false, err
	}

	err = transaction.UpdateReceivedMessage(receivedMessage)
	if err != nil {
		transaction.Dispose()
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
