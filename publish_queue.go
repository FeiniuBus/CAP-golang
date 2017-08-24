package cap

type PublishQueuer struct{
	StorageConnection IStorageConnection
}

func NewPublishQueuer(storageConnection IStorageConnection) *PublishQueuer{
	queuer := &PublishQueuer{StorageConnection:storageConnection}
	return queuer
}

func (this *PublishQueuer) Execute() error{
	var message *CapPublishedMessage
	var err error
	for{
		transaction,err := this.StorageConnection.CreateTransaction()
		message,err = this.StorageConnection.GetNextPublishedMessageToBeEnqueued()
		if err != nil || message == nil || message.Id==0{
			break
		}
		err = transaction.EnqueuePublishedMessage(message)
		if err != nil {
			break
		}
		message.StatusName = "Enqueued"
		err = transaction.UpdatePublishedMessage(message)
		if err != nil {
			break
		}
		err = transaction.Commit()
		if err != nil {
			break
		}
	}
	if err != nil{
		return err
	}
	return nil
}