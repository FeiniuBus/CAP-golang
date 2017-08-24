package cap

type IState interface {
	GetExpiresAfter() int32
	GetName() string
	ApplyReceivedMessage(message *CapReceivedMessage, transaction IStorageTransaction) error
	ApplyPublishedMessage(message *CapPublishedMessage, transaction IStorageTransaction) error
}