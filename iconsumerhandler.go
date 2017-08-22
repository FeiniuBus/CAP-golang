package cap

type IConsumerHandler interface {
	Start()
	Close()
}