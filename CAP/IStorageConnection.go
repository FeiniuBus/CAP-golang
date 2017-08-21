package cap

type IStorageConnection interface{

	CreateTransaction() IStorageTransaction;

	FetchNextMessage(dbConnection interface{}) (IFetchedMessage,error);

	GetFailedPublishedMessages() ([]*CapPublishedMessage,error);

	GetNextPublishedMessageToBeEnqueued() (*CapPublishedMessage,error);

	GetNextReceviedMessageToBeEnqueued() (*CapReceivedMessage,error);

	GetPublishedMessage(id int) (*CapPublishedMessage, error);

	GetReceivedMessage(id int) (*CapReceivedMessage, error);

	StoreReceivedMessage(message *CapReceivedMessage) error;
	
}