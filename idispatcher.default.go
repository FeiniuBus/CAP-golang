package cap

type DefaultDispatcher struct {
	StorageConnectionFactory *StorageConnectionFactory
	CapOptions               *CapOptions
	QueueExecutorFactory     *QueueExecutorFactory
}

func NewDefaultDispatcher(capOptions *CapOptions,
	storageConnectionFactory *StorageConnectionFactory,
) IProcessor {
	return &DefaultDispatcher{
		StorageConnectionFactory: storageConnectionFactory,
		CapOptions:               capOptions,
	}
}

func (this *DefaultDispatcher) Process(context *ProcessingContext) (*ProcessResult, error) {
	err := context.ThrowIfStopping()
	if err != nil {
		return nil, err
	}

	err = this.ProcessCore(context)
	if err != nil {
		return nil, err
	}

	return ProcessSleeping(this.CapOptions.PoolingDelay), nil
}

func (this *DefaultDispatcher) ProcessCore(context *ProcessingContext) error {
	_, err := this.step(context)
	return err
}

func (this *DefaultDispatcher) step(context *ProcessingContext) (bool, error) {
	for {
		conn, err := this.StorageConnectionFactory.CreateStorageConnection(this.CapOptions)
		if err != nil {
			return false, nil
		}

		fetched, err := conn.FetchNextMessage()
		if err != nil {
			return false, nil
		}

		var messageType string
		if fetched.GetMessageType() == 0 {
			messageType = PUBLISH
		} else {
			messageType = SUBSCRIBE
		}

		queueExecutor := this.QueueExecutorFactory.GetInstance(messageType)
		err = queueExecutor.Execute(conn, fetched)
		if err != nil {
			return false, err
		}

		err = fetched.Dispose()

		return true, err
	}
}
