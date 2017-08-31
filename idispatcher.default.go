package cap

// DefaultDispatcher ...
type DefaultDispatcher struct {
	StorageConnectionFactory *StorageConnectionFactory
	CapOptions               *CapOptions
	QueueExecutorFactory     IQueueExecutorFactory
}

// NewDefaultDispatcher ...
func NewDefaultDispatcher(capOptions *CapOptions,
	storageConnectionFactory *StorageConnectionFactory,
	factory IQueueExecutorFactory,
) IProcessor {
	return &DefaultDispatcher{
		StorageConnectionFactory: storageConnectionFactory,
		CapOptions:               capOptions,
		QueueExecutorFactory:     factory,
	}
}

// Process ...
func (dispatcher *DefaultDispatcher) Process(context *ProcessingContext) (*ProcessResult, error) {
	for {
		err := context.ThrowIfStopping()
		if err != nil {
			return nil, err
		}

		worked, err := dispatcher.ProcessCore(context)
		if err != nil {
			return nil, err
		}

		if !worked {
			break
		}
	}

	return ProcessSleeping(dispatcher.CapOptions.PoolingDelay), nil
}

// ProcessCore ...
func (dispatcher *DefaultDispatcher) ProcessCore(context *ProcessingContext) (bool, error) {
	return dispatcher.step(context)
}

func (dispatcher *DefaultDispatcher) step(context *ProcessingContext) (bool, error) {

	conn, err := dispatcher.StorageConnectionFactory.CreateStorageConnection(dispatcher.CapOptions)
	if err != nil {
		return false, err
	}

	fetched, err := conn.FetchNextMessage()
	if err != nil {
		return false, err
	}
	defer fetched.Dispose()

	if fetched.GetMessageId() == 0 {
		return false, nil
	}

	var messageType string
	if fetched.GetMessageType() == 0 {
		messageType = PUBLISH
	} else {
		messageType = SUBSCRIBE
	}

	queueExecutor := dispatcher.QueueExecutorFactory.GetInstance(messageType)
	err = queueExecutor.Execute(conn, fetched)
	if err != nil {
		return false, err
	}

	return true, err
}
