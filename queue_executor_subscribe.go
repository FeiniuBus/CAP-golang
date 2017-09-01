package cap

// QueueExecutorSubscribe ...
type QueueExecutorSubscribe struct {
	IQueueExecutor
	Register *CallbackRegister
	logger   ILogger
}

// NewQueueExecutorSubscribe ..
func NewQueueExecutorSubscribe(register *CallbackRegister) *QueueExecutorSubscribe {
	executor := &QueueExecutorSubscribe{
		Register: register,
	}
	executor.logger = GetLoggerFactory().CreateLogger(executor)
	return executor
}

// Execute ...
func (executor *QueueExecutorSubscribe) Execute(connection IStorageConnection, feched IFetchedMessage) error {
	message, err := connection.GetReceivedMessage(feched.GetMessageId())
	if err != nil {
		executor.logger.LogData(LevelError, "[Execute]"+err.Error(), message)
		return err
	}

	if message == nil || message.Id == 0 {
		return nil
	}

	stateChanger := NewStateChanger()
	transaction, err := connection.CreateTransaction()
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	defer transaction.Dispose()

	err = stateChanger.ChangeReceivedMessageState(message, NewProcessingState(), transaction)
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	var newState IState
	err = executor.executeSubscribeAsync(message)
	if err != nil {
		shouldRetry, err := executor.updateMessageForRetryAsync(message, connection)
		if err != nil {
			executor.logger.Log(LevelError, "[Execute]"+err.Error())
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
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}
	defer transaction.Dispose()

	err = stateChanger.ChangeReceivedMessageState(message, newState, transaction)
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	err = feched.RemoveFromQueue()
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	return nil
}

func (executor *QueueExecutorSubscribe) executeSubscribeAsync(receivedMessage *CapReceivedMessage) error {
	cb, err := executor.Register.Get(receivedMessage.Group, receivedMessage.Name)
	if err != nil {
		executor.logger.Log(LevelError, "[executeSubscribeAsync]"+err.Error())
		return nil
	}
	return cb.Handle(receivedMessage)
}

func (executor *QueueExecutorSubscribe) updateMessageForRetryAsync(receivedMessage *CapReceivedMessage, conn IStorageConnection) (bool, error) {
	retryBehavior := DefaultRetry
	receivedMessage.Retries = receivedMessage.Retries + 1
	if receivedMessage.Retries >= int(retryBehavior.RetryCount) {
		return false, nil
	}

	due := receivedMessage.Added + int64(retryBehavior.RetryIn(int32(receivedMessage.Retries)))
	receivedMessage.ExpiresAt = int(due)

	transaction, err := conn.CreateTransaction()
	if err != nil {
		executor.logger.Log(LevelError, "[updateMessageForRetryAsync]"+err.Error())
		return false, err
	}

	err = transaction.UpdateReceivedMessage(receivedMessage)
	if err != nil {
		transaction.Dispose()
		executor.logger.Log(LevelError, "[updateMessageForRetryAsync]"+err.Error())
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		executor.logger.Log(LevelError, "[updateMessageForRetryAsync]"+err.Error())
		return false, err
	}

	return true, nil
}
