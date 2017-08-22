package cap

type IStorageTransaction interface{

	Commit() error;

	EnqueuePublishedMessage(message *CapPublishedMessage) error;

	EnqueueReceivedMessage(message *CapReceivedMessage) error;

	UpdatePublishedMessage(message *CapPublishedMessage) error;

	UpdateReceivedMessage(message *CapReceivedMessage) error;
}