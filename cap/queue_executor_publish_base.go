package cap

// QueueExecutorPublish ...
type QueueExecutorPublish struct {
	IQueueExecutor
	StateChanger    IStateChanger
	PublishDelegate IPublishDelegate
	logger          ILogger
}

// NewQueueExecutorPublish ...
func NewQueueExecutorPublish(stateChanger IStateChanger, delegate IPublishDelegate) *QueueExecutorPublish {
	executor := &QueueExecutorPublish{
		StateChanger:    stateChanger,
		PublishDelegate: delegate,
	}
	executor.logger = GetLoggerFactory().CreateLogger(executor)
	return executor
}

// Execute ...
func (executor *QueueExecutorPublish) Execute(connection IStorageConnection, feched IFetchedMessage) error {
	message, err := connection.GetPublishedMessage(feched.GetMessageId())
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	if message == nil || message.Id == 0 {
		return nil
	}

	transaction, err := connection.CreateTransaction()
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}
	defer transaction.Dispose()

	err = executor.StateChanger.ChangePublishedMessage(message, NewProcessingState(), transaction)
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		return err
	}

	err = executor.PublishDelegate.Publish(message.Name, message.Content)

	var newState IState = nil
	if err != nil {
		executor.logger.Log(LevelError, "[Execute]"+err.Error())
		shouldRetry, err := executor.UpdateMessageForRetry(message, connection)
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

	err = executor.StateChanger.ChangePublishedMessage(message, newState, transaction)
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

// UpdateMessageForRetry ...
func (executor *QueueExecutorPublish) UpdateMessageForRetry(message *CapPublishedMessage, connection IStorageConnection) (bool, error) {
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
