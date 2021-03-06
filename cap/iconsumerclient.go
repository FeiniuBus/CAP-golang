package cap

type ErrorHanlder func (str string)
type ReceiveHanlder func (ctx MessageContext)

type IConsumerClient interface {
	Subscribe(topics []string)
	Listening(timeoutSecs int, ch chan bool)
	Commit(context MessageContext)
	Close()
	SetOnReceive(onReceive ReceiveHanlder)
	SetOnError(onError ErrorHanlder)
}