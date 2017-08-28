package cap

type DefaultDispatcher struct {
	StorageConnectionFactory *StorageConnectionFactory
	CapOptions               *CapOptions
	QueueExecutorFactory     *QueueExecutorFactory
}

func NewDefaultDispatcher(storageConnectionFactory *StorageConnectionFactory,
	capOptions *CapOptions) *DefaultDispatcher {
	return &DefaultDispatcher{
		StorageConnectionFactory: storageConnectionFactory,
		CapOptions:               capOptions,
	}
}

func (this *DefaultDispatcher) Process(context *ProcessingContext) error {
	err := context.ThrowIfStopping()
	if err != nil {
		return err
	}

	return this.ProcessCore(context)
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
