package cap

// FailedJobProcessor ...
type FailedJobProcessor struct {
	Options                  *CapOptions
	StateChanger             IStateChanger
	StorageConnectionFactory *StorageConnectionFactory
	logger                   ILogger
}

// NewFailedJobProcessor...
func NewFailedJobProcessor(capOptions *CapOptions, storageConnectionFactory *StorageConnectionFactory) IProcessor {
	processor := &FailedJobProcessor{
		Options:                  capOptions,
		StorageConnectionFactory: storageConnectionFactory,
		StateChanger:             NewStateChanger(),
	}
	processor.logger = GetLoggerFactory().CreateLogger(processor)
	return processor
}

// Process ...
func (processor *FailedJobProcessor) Process(context *ProcessingContext) (*ProcessResult, error) {
	connection, err := processor.StorageConnectionFactory.CreateStorageConnection(processor.Options)
	if err != nil {
		processor.logger.Log(LevelError, "[Process]"+err.Error())
		return nil, err
	}

	err = processor.ProcessPublishedMessage(connection)
	if err != nil {
		processor.logger.Log(LevelError, "[Process]"+err.Error())
		return nil, err
	}
	err = processor.ProcessReceivedMessage(connection)
	if err != nil {
		processor.logger.Log(LevelError, "[Process]"+err.Error())
		return nil, err
	}

	return ProcessSleeping(processor.Options.PoolingDelay), nil
}

// ProcessPublishedMessage ...
func (processor *FailedJobProcessor) ProcessPublishedMessage(connection IStorageConnection) error {
	failedPublishedMessages, err := connection.GetFailedLockedMessages(0)
	if err != nil {
		processor.logger.Log(LevelError, "[ProcessPublishedMessage]"+err.Error())
		return err
	}
	defer failedPublishedMessages.Dispose()
	err = failedPublishedMessages.ChangeStates(NewEnqueuedState())
	if err != nil {
		failedPublishedMessages.Rollback()
		processor.logger.Log(LevelError, "[ProcessPublishedMessage]"+err.Error())
		return err
	}
	err = failedPublishedMessages.Commit()
	if err != nil {
		failedPublishedMessages.Rollback()
		processor.logger.Log(LevelError, "[ProcessPublishedMessage]"+err.Error())
		return err
	}
	return nil
}

// ProcessReceivedMessage ...
func (processor *FailedJobProcessor) ProcessReceivedMessage(connection IStorageConnection) error {
	failedReceivedMessages, err := connection.GetFailedLockedMessages(1)
	if err != nil {
		processor.logger.Log(LevelError, "[ProcessReceivedMessage]"+err.Error())
		return err
	}
	defer failedReceivedMessages.Dispose()
	err = failedReceivedMessages.ChangeStates(NewEnqueuedState())
	if err != nil {
		failedReceivedMessages.Rollback()
		processor.logger.Log(LevelError, "[ProcessReceivedMessage]"+err.Error())
		return err
	}
	err = failedReceivedMessages.Commit()
	if err != nil {
		failedReceivedMessages.Rollback()
		processor.logger.Log(LevelError, "[ProcessReceivedMessage]"+err.Error())
		return err
	}
	return nil
}
