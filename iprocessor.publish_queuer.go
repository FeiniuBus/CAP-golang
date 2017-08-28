package cap

type PublishQueuer struct{
	StateChanger IStateChanger
	Options *CapOptions
	StorageConnectionFactory *StorageConnectionFactory
}

func (this *PublishQueuer) Process(context *ProcessingContext) error{
	var message *CapPublishedMessage
	connection,err := this.StorageConnectionFactory.CreateStorageConnection(this.Options)
	if err != nil {
		return err
	}
	
	looper := NewLooper()
	err = looper.While(func()bool{
		message, err = connection.GetNextPublishedMessageToBeEnqueued()
		if err != nil{
			return false
		}
		return (context.IsStopping == false)
	},func()error{
		state := NewScheduledState()
		transaction,err := connection.CreateTransaction()
		if err != nil {
			return err
		}
		err = this.StateChanger.ChangePublishedMessage(message,state,transaction)
		if err != nil {
			return err
		}
		err = transaction.Commit()
		if err != nil {
			return err
		}
	})
	if err != nil {
		return err
	}
	err = context.ThrowIfStopping()
	if err != nil {
		return err
	}
	return nil
}