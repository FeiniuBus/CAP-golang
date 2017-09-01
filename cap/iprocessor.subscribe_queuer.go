package cap

// SubscribeQueuer ...
type SubscribeQueuer struct {
	Options                  *CapOptions
	StorageConnectionFactory *StorageConnectionFactory
	logger                   ILogger
}

// NewSubscribeQueuer ...
func NewSubscribeQueuer(options *CapOptions, connectionFactory *StorageConnectionFactory) IProcessor {
	queuer := &SubscribeQueuer{
		Options:                  options,
		StorageConnectionFactory: connectionFactory,
	}
	queuer.logger = GetLoggerFactory().CreateLogger(queuer)
	return queuer
}

// Process ...
func (processor *SubscribeQueuer) Process(context *ProcessingContext) (*ProcessResult, error) {
	conn, err := processor.StorageConnectionFactory.CreateStorageConnection(processor.Options)
	if err != nil {
		processor.logger.Log(LevelError, "[Process]"+err.Error())
		return nil, err
	}

	for {
		if context.IsStopping {
			break
		}
		//message, err := conn.GetNextReceviedMessageToBeEnqueued()
		message, err := conn.GetNextLockedMessageToBeEnqueued(1)
		if err != nil {
			processor.logger.Log(LevelError, "[Process]"+err.Error())
			return nil, err
		}
		if message == nil || message.GetMessage().(*CapReceivedMessage).Id == 0 {
			processor.logger.Log(LevelInfomation, "[Process]Message is nil, task canceled.")
			break
		}

		//transaction, err := conn.CreateTransaction()
		// if err != nil {
		// 	processor.logger.Log(LevelError, "[Process]"+err.Error())
		// 	return nil, err
		// }

		//stateChanger := NewStateChanger()
		//err = stateChanger.ChangeReceivedMessageState(message, NewEnqueuedState(), transaction)
		err = message.ChangeState(NewEnqueuedState())
		if err != nil {
			//transaction.Dispose()
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
		//transaction.Dispose()
		message.Dispose()
	}
	err = context.ThrowIfStopping()
	if err != nil {
		processor.logger.Log(LevelError, "[Process]"+err.Error())
		return nil, err
	}
	return ProcessSleeping(processor.Options.PoolingDelay), nil
}
