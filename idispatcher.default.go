package cap

// DefaultDispatcher ...
type DefaultDispatcher struct {
	StorageConnectionFactory *StorageConnectionFactory
	CapOptions               *CapOptions
	QueueExecutorFactory     *QueueExecutorFactory
}

// NewDefaultDispatcher ...
func NewDefaultDispatcher(capOptions *CapOptions,
	storageConnectionFactory *StorageConnectionFactory,
) IProcessor {
	return &DefaultDispatcher{
		StorageConnectionFactory: storageConnectionFactory,
		CapOptions:               capOptions,
	}
}

// Process ...
func (dispatcher *DefaultDispatcher) Process(context *ProcessingContext) (*ProcessResult, error) {
	err := context.ThrowIfStopping()
	if err != nil {
		return nil, err
	}

	err = dispatcher.ProcessCore(context)
	if err != nil {
		return nil, err
	}

	return ProcessSleeping(dispatcher.CapOptions.PoolingDelay), nil
}

// ProcessCore ...
func (dispatcher *DefaultDispatcher) ProcessCore(context *ProcessingContext) error {
	_, err := dispatcher.step(context)
	return err
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

	err = fetched.Dispose()

	return true, err
}
