package cap


type FailedJobProcessor struct{
	Options *CapOptions
	StateChanger IStateChanger
	StorageConnectionFactory *StorageConnectionFactory
}

func (this *FailedJobProcessor) Process(context *ProcessingContext) error{
	connection, err := this.StorageConnectionFactory.CreateStorageConnection(this.Options)
	if err != nil{
		return err
	}
	err = this.ProcessPublishedMessage(connection)
	if err != nil{
		return err
	}
	err = this.ProcessReceivedMessage(connection)
	if err != nil{
		return err
	}
	return nil
}

func (this *FailedJobProcessor) ProcessPublishedMessage(connection IStorageConnection) error{
	hasException := false
	messages, err := connection.GetFailedPublishedMessages()
	if err != nil{
		return err
	}
	length := len(messages)
	for i:=0;i<length;i++ {
		message := messages[i]
		if hasException == false {
			//TODO: failed callback
		}
		transaction, err := connection.CreateTransaction()
		if err != nil {
			return err
		}
		err = this.StateChanger.ChangePublishedMessage(message,NewEnqueuedState(),transaction)
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

func (this *FailedJobProcessor) ProcessReceivedMessage(connection IStorageConnection) error{
	hasException := false
	messages, err := connection.GetFailedReceivedMessages()
	if err != nil{
		return err
	}
	length := len(messages)
	for i:=0;i<length;i++ {
		message := messages[i]
		if hasException == false {
			//TODO: failed callback
		}
		transaction, err := connection.CreateTransaction()
		if err != nil {
			return err
		}
		err = this.StateChanger.ChangeReceivedMessageState(message,NewEnqueuedState(),transaction)
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