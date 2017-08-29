package rabbitmq

import (
	"github.com/FeiniuBus/capgo"
)

func Prepare(bootstrapper *cap.Bootstrapper, rabbitMQOptions RabbitMQOptions) {
	consumerHandler := NewConsumerHandlerRabbitMQ(
		rabbitMQOptions,
		bootstrapper.CapOptions,
		bootstrapper.Register,
		bootstrapper.ConnectionFactory,
	)
	bootstrapper.Servers = append(bootstrapper.Servers, consumerHandler)
	bootstrapper.QueueExecutorFactory.SetPublishQueueExecutorCreateDelegate(func() cap.IQueueExecutor {
		return cap.NewQueueExecutorPublish(cap.NewStateChanger(), NewPublishQueueExecutor(cap.NewStateChanger(), &rabbitMQOptions))
	})
}
