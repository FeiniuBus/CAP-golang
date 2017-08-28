package cap

type IStorageConnection interface{

	CreateTransaction() (IStorageTransaction,error);

	FetchNextMessage() (IFetchedMessage,error);

	GetFailedPublishedMessages() ([]*CapPublishedMessage,error);

	GetFailedReceivedMessages() ([]*CapReceivedMessage,error);

	GetNextPublishedMessageToBeEnqueued() (*CapPublishedMessage,error);

	GetNextReceviedMessageToBeEnqueued() (*CapReceivedMessage,error);

	GetPublishedMessage(id int) (*CapPublishedMessage, error);

	GetReceivedMessage(id int) (*CapReceivedMessage, error);

	StoreReceivedMessage(message *CapReceivedMessage) error;
	
}