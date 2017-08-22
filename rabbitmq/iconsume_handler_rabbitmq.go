package rabbitmq

import (
	"runtime"
	"github.com/FeiniuBus/capgo"	
)

type ConsumerHandlerRabbitMQ struct {
	cap.IConsumerHandler

	ConsumerClientFactory cap.IConsumerClientFactory
	RabbitOptions RabbitMQOptions
}

func NewConsumerHandlerRabbitMQ(rabbitOptions RabbitMQOptions) *ConsumerHandlerRabbitMQ {
	clientFactory := NewRabbitConsumeClientFactory(rabbitOptions)
	
	rtv := &ConsumerHandlerRabbitMQ {
		RabbitOptions: rabbitOptions ,
		ConsumerClientFactory: clientFactory,
	}

	runtime.SetFinalizer(rtv, rtv.Close)

	return rtv 
}

func (this *ConsumerHandlerRabbitMQ) Start() {

}

func (this *ConsumerHandlerRabbitMQ) Close() {
	
}

func (this *ConsumerHandlerRabbitMQ) Pulse() {
	
}

