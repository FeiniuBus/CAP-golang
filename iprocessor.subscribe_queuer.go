package cap

type SubscribeQueuer struct {
	Options                  *CapOptions
	StorageConnectionFactory *StorageConnectionFactory
}

func NewSubscribeQueuer(options *CapOptions, connectionFactory *StorageConnectionFactory) *SubscribeQueuer {
	return &SubscribeQueuer{
		Options:                  options,
		StorageConnectionFactory: connectionFactory,
	}
}

func (this *SubscribeQueuer) Process(context *ProcessingContext) error {
	for {
		if context.IsStopping {
			return nil
		}

		conn, err := this.StorageConnectionFactory.CreateStorageConnection(this.Options)
		if err != nil {
			return err
		}

		message, err := conn.GetNextReceviedMessageToBeEnqueued()
		if err != nil {
			return err
		}
		if message == nil {
			return nil
		}

		transaction, err := conn.CreateTransaction()
		if err != nil {
			return err
		}

		stateChanger := NewStateChanger()
		err = stateChanger.ChangeReceivedMessageState(message, NewEnqueuedState(), transaction)
		if err != nil {
			return err
		}

		err = transaction.Commit()
		if err != nil {
			return err
		}

		err = context.ThrowIfStopping()
		if err != nil {
			return err
		}
	}
}
