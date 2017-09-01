package cap

type IReceivedMessageHandler interface{
	Handle(message interface{})
}