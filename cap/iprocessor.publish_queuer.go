package cap

// PublishQueuer blablabla .
type PublishQueuer struct {
	StateChanger             IStateChanger
	Options                  *CapOptions
	StorageConnectionFactory *StorageConnectionFactory
	logger                   ILogger
}

// NewPublishQueuer bla.
func NewPublishQueuer(capOptions *CapOptions, storageConnectionFactory *StorageConnectionFactory) IProcessor {
	publisher := &PublishQueuer{
		Options:                  capOptions,
		StorageConnectionFactory: storageConnectionFactory,
		StateChanger:             NewStateChanger(),
	}
	publisher.logger = GetLoggerFactory().CreateLogger(publisher)
	return publisher
}

// Process blablabla.
func (processor *PublishQueuer) Process(context *ProcessingContext) (*ProcessResult, error) {
	connection, err := processor.StorageConnectionFactory.CreateStorageConnection(processor.Options)
	if err != nil {
		return nil, err
	}

	for {
		if context.IsStopping {
			break
		}
		message, err := connection.GetNextLockedMessageToBeEnqueued(0)
		if err != nil {
			processor.logger.Log(LevelError, "[Process]"+err.Error())
			return nil, err
		}

		if message == nil || message.GetMessage().(*CapPublishedMessage).Id == 0 {
			err = context.ThrowIfStopping()
			if err != nil {
				processor.logger.Log(LevelError, "[Process]"+err.Error())
				return nil, err
			}
			break
		}

		//state := NewEnqueuedState()
		//transaction, err := connection.CreateTransaction()
		// if err != nil {
		// 	processor.logger.Log(LevelError, "[Process]"+err.Error())
		// 	return nil, err
		// }
		// defer transaction.Dispose()

		//err = processor.StateChanger.ChangePublishedMessage(message, state, transaction)
		err = message.ChangeState(NewEnqueuedState())
		if err != nil {
			processor.logger.Log(LevelError, "[Process]"+err.Error())
			message.Rollback()
			message.Dispose()
			return nil, err
		}

		//err = transaction.Commit()
		err = message.Commit()
		if err != nil {
			processor.logger.Log(LevelError, "[Process]"+err.Error())
			message.Dispose()
			return nil, err
		}

		err = context.ThrowIfStopping()
		if err != nil {
			processor.logger.Log(LevelError, "[Process]"+err.Error())
			message.Dispose()
			return nil, err
		}

		message.Dispose()
	}
	return ProcessSleeping(processor.Options.PoolingDelay), nil
}
