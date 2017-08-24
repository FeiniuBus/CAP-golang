package rabbitmq

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/FeiniuBus/capgo"	
)

type RabbitMQConsumerClient struct {
	cap.IConsumerClient

	QueueName 		string
	ConnectString 	string	
	Options 		*RabbitMQOptions
	Connection 		*amqp.Connection
	Channel			*amqp.Channel

	OnReceive		cap.ReceiveHanlder
	OnError			cap.ErrorHanlder
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(err)
	}
}

func NewClient(queueName string , options *RabbitMQOptions) *RabbitMQConsumerClient {
	rtv := &RabbitMQConsumerClient {
		Options: options ,
		QueueName: queueName ,
		ConnectString: ConnectString(options) ,
	}

	rtv.InitClient()

	return rtv
}

func (this *RabbitMQConsumerClient) InitClient() {
	conn, err := amqp.Dial(this.ConnectString)
	failOnError(err, "Fail to connect to rabbit mq " + this.ConnectString)
	this.Connection = conn

	ch, err := conn.Channel()
	failOnError(err, "Fail to create channel " + this.ConnectString)
	this.Channel = ch

	err = this.Channel.ExchangeDeclare(
		this.Options.TopicExchangeName , 
		this.Options.ExchangeType , 
		true ,
		false , 
		false ,
		false ,
		nil, 
	)

	failOnError(err, "Fail to declare exchange")

	args := amqp.Table {
		"x-message-ttl": this.Options.QueueMessageExpires ,
	}

	_, err = this.Channel.QueueDeclare(
		this.QueueName ,
		true ,
		false ,
		false ,
		false ,
		args ,
	)

	failOnError(err, "fail to declare queue")
}

func (this *RabbitMQConsumerClient) Close() {
	if this.Connection != nil {
		this.Connection.Close()
		this.Connection = nil
	}

	if this.Channel != nil {
		this.Channel.Close()
		this.Channel = nil
	}
}

func (this *RabbitMQConsumerClient) Subscribe(topics []string) {
	for _, value := range topics {
		this.Channel.QueueBind(
			this.QueueName , 
			value ,
			this.Options.TopicExchangeName ,
			false ,
			nil)
	}
}

func (this *RabbitMQConsumerClient) Listening(timeoutSecs int, done chan bool) {
	msgs, err := this.Channel.Consume(
		this.QueueName ,
		"" ,
		false ,
		false ,
		false ,
		false ,
		nil )

	failOnError(err, "Consume fail")

	handleReceive(this, msgs, done)
}

func handleReceive(client *RabbitMQConsumerClient, deliveries <-chan amqp.Delivery, done <-chan bool) {
	go func ( deliveries <-chan amqp.Delivery, done <-chan bool){
		for {
			select {
				case delivery := <-deliveries:
					context := cap.MessageContext {
						Group: client.QueueName,
						Name: delivery.RoutingKey,
						Content: string(delivery.Body),
						Tag: delivery.DeliveryTag,
					}

					if client.OnReceive != nil {
						client.OnReceive(context)
					}

				case <-done:
					return
			}
		}
	}(deliveries, done)
}

func (this *RabbitMQConsumerClient) Commit(context cap.MessageContext) {
	this.Channel.Ack(context.Tag, false)
}

func (this *RabbitMQConsumerClient) SetOnReceive(onReceive cap.ReceiveHanlder) {
	this.OnReceive = onReceive
}

func (this *RabbitMQConsumerClient) SetOnError(onError cap.ErrorHanlder) {
	this.OnError = onError
}