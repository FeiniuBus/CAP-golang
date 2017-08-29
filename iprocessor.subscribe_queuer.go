package cap

// SubscribeQueuer ...
type SubscribeQueuer struct {
	Options                  *CapOptions
	StorageConnectionFactory *StorageConnectionFactory
}

// NewSubscribeQueuer ...
func NewSubscribeQueuer(options *CapOptions, connectionFactory *StorageConnectionFactory) IProcessor {
	return &SubscribeQueuer{
		Options:                  options,
		StorageConnectionFactory: connectionFactory,
	}
}

// Process ...
func (processor *SubscribeQueuer) Process(context *ProcessingContext) (*ProcessResult, error) {
	conn, err := processor.StorageConnectionFactory.CreateStorageConnection(processor.Options)
	if err != nil {
		return nil, err
	}

	message, err := conn.GetNextReceviedMessageToBeEnqueued()
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, nil
	}

	transaction, err := conn.CreateTransaction()
	if err != nil {
		return nil, err
	}

	stateChanger := NewStateChanger()
	err = stateChanger.ChangeReceivedMessageState(message, NewEnqueuedState(), transaction)
	if err != nil {
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	err = context.ThrowIfStopping()
	if err != nil {
		return nil, err
	}

	return ProcessSleeping(processor.Options.PoolingDelay), nil
}
