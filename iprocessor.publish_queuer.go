package cap

// PublishQueuer blablabla .
type PublishQueuer struct {
	StateChanger             IStateChanger
	Options                  *CapOptions
	StorageConnectionFactory *StorageConnectionFactory
}

// NewPublishQueuer bla.
func NewPublishQueuer(capOptions *CapOptions, storageConnectionFactory *StorageConnectionFactory) IProcessor {
	publisher := &PublishQueuer{
		Options:                  capOptions,
		StorageConnectionFactory: storageConnectionFactory,
		StateChanger:             NewStateChanger(),
	}
	return publisher
}

// Process blablabla.
func (processor *PublishQueuer) Process(context *ProcessingContext) (*ProcessResult, error) {
	var message *CapPublishedMessage
	connection, err := processor.StorageConnectionFactory.CreateStorageConnection(processor.Options)

	message, err = connection.GetNextPublishedMessageToBeEnqueued()
	if err != nil {
		return nil, err
	}

	if message == nil || message.Id == 0 {
		err = context.ThrowIfStopping()
		if err != nil {
			return nil, err
		}

		return ProcessSleeping(processor.Options.PoolingDelay), nil
	}

	state := NewEnqueuedState()
	transaction, err := connection.CreateTransaction()
	if err != nil {
		return nil, err
	}

	err = processor.StateChanger.ChangePublishedMessage(message, state, transaction)
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
