package cap

// FailedJobProcessor ...
type FailedJobProcessor struct {
	Options                  *CapOptions
	StateChanger             IStateChanger
	StorageConnectionFactory *StorageConnectionFactory
}

// NewFailedJobProcessor...
func NewFailedJobProcessor(capOptions *CapOptions, storageConnectionFactory *StorageConnectionFactory) IProcessor {
	return &FailedJobProcessor{
		Options:                  capOptions,
		StorageConnectionFactory: storageConnectionFactory,
		StateChanger:             NewStateChanger(),
	}
}

// Process ...
func (processor *FailedJobProcessor) Process(context *ProcessingContext) (*ProcessResult, error) {
	connection, err := processor.StorageConnectionFactory.CreateStorageConnection(processor.Options)
	if err != nil {
		return nil, err
	}
	err = processor.ProcessPublishedMessage(connection)
	if err != nil {
		return nil, err
	}
	err = processor.ProcessReceivedMessage(connection)
	if err != nil {
		return nil, err
	}

	return ProcessSleeping(processor.Options.PoolingDelay), nil
}

// ProcessPublishedMessage ...
func (processor *FailedJobProcessor) ProcessPublishedMessage(connection IStorageConnection) error {
	hasException := false
	messages, err := connection.GetFailedPublishedMessages()
	if err != nil {
		return err
	}
	length := len(messages)
	for i := 0; i < length; i++ {
		message := messages[i]
		if hasException == false {
			//TODO: failed callback
		}
		transaction, err := connection.CreateTransaction()
		if err != nil {
			return err
		}
		err = processor.StateChanger.ChangePublishedMessage(message, NewEnqueuedState(), transaction)
		if err != nil {
			return err
		}
		err = transaction.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessReceivedMessage ...
func (processor *FailedJobProcessor) ProcessReceivedMessage(connection IStorageConnection) error {
	hasException := false
	messages, err := connection.GetFailedReceivedMessages()
	if err != nil {
		return err
	}
	length := len(messages)
	for i := 0; i < length; i++ {
		message := messages[i]
		if hasException == false {
			//TODO: failed callback
		}
		transaction, err := connection.CreateTransaction()
		if err != nil {
			return err
		}
		err = processor.StateChanger.ChangeReceivedMessageState(message, NewEnqueuedState(), transaction)
		if err != nil {
			return err
		}
		err = transaction.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}
