package cap

type IFetchedMessage interface{
	GetMessageId()(messageId int)

	GetMessageType()(messageType int)

	RemoveFromQueue() error
	Requeue() error

	Dispose()
}