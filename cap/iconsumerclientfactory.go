package cap

type IConsumerClientFactory interface {
	Create(group string) IConsumerClient
}