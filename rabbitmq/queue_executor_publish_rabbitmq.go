package rabbitmq

import (
	"github.com/FeiniuBus/capgo"
	"github.com/streadway/amqp"
)

type PublishQueueExecutor struct {
	cap.QueueExecutorPublishBase
	RabbitOptions *RabbitMQOptions
}

func NewPublishQueueExecutor(stateChanger cap.IStateChanger, rabbitOptions *RabbitMQOptions) *PublishQueueExecutor {
	rtv := &PublishQueueExecutor{
		RabbitOptions: rabbitOptions,
	}
	rtv.StateChanger = stateChanger

	return rtv
}

func (this *PublishQueueExecutor) Execute(connection cap.IStorageConnection, feched cap.IFetchedMessage) error {
	return this.QueueExecutorPublishBase.Execute(connection, feched)
}

func (this *PublishQueueExecutor) Publish(keyName, content string) error {
	connectString := ConnectString(this.RabbitOptions)

	conn, err := amqp.Dial(connectString)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		this.RabbitOptions.TopicExchangeName,
		this.RabbitOptions.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		this.RabbitOptions.TopicExchangeName,
		keyName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
