package cap_rabbitmq

import (
	"../cap"
)

type RabbitMQConsumerClientFactory struct {
	cap.IConsumerClientFactory
	Options RabbitMQOptions
}

func NewRabbitConsumeClientFactory(options RabbitMQOptions) RabbitMQConsumerClientFactory {
	return RabbitMQConsumerClientFactory{
		Options: options ,
	}
}

func (this *RabbitMQConsumerClientFactory) Create(group string) cap.IConsumerClient {
	return NewClient(group, &this.Options)
}
