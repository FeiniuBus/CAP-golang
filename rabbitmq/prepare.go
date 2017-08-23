package rabbitmq

import (
	"github.com/FeiniuBus/capgo"
)

func Prepare(bootstrapper *cap.Bootstrapper, rabbitMQOptions RabbitMQOptions ) {
	consumerHandler := NewConsumerHandlerRabbitMQ(
		rabbitMQOptions,
		bootstrapper.CapOptions,
		bootstrapper.Register,
		bootstrapper.ConnectionFactory,
	)
	bootstrapper.Servers = append(bootstrapper.Servers, consumerHandler)
}