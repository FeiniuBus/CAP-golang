package rabbitmq

import (
	"runtime"
	"github.com/FeiniuBus/capgo"	
)

const (
	DefaultPollingDelay = 60
)

type ConsumerHandlerRabbitMQ struct {
	cap.IConsumerHandler

	ConsumerClientFactory 	cap.IConsumerClientFactory
	RabbitOptions 			RabbitMQOptions
	Register 				*cap.CallbackRegister
	PollingDelay 			uint32
	Done					chan bool
	Clients					[]cap.IConsumerClient
	ConnectionFactory		*cap.StorageConnectionFactory
	CapOptions				cap.CapOptions
}

func NewConsumerHandlerRabbitMQ(
	rabbitOptions RabbitMQOptions,
	capOptions cap.CapOptions, 
	register *cap.CallbackRegister, 
	connectionFactory *cap.StorageConnectionFactory) *ConsumerHandlerRabbitMQ {
	clientFactory := NewRabbitConsumeClientFactory(rabbitOptions)
	
	rtv := &ConsumerHandlerRabbitMQ {
		RabbitOptions: rabbitOptions ,
		ConsumerClientFactory: clientFactory ,
		Register: register ,
		PollingDelay: DefaultPollingDelay,
		Done: make(chan bool) ,
		Clients: make([]cap.IConsumerClient, 0) ,
		ConnectionFactory: connectionFactory,
		CapOptions: capOptions ,
	}

	return rtv 
}

func (this *ConsumerHandlerRabbitMQ) Start() {
	for group, groupItems := range this.Register.Routers {
		client := this.ConsumerClientFactory.Create(group)
		
		names := []string{}
		for name, _ := range groupItems {
			names = append(names, name)
		}
		
		this.registerMessageProcessor(client)
		client.Subscribe(names)
		client.Listening(60, this.Done)

		this.Clients = append(this.Clients, client)
	}
}

func (this *ConsumerHandlerRabbitMQ) Close() {
	for _, client := range this.Clients {
		client.Close()
	}
}

func (this *ConsumerHandlerRabbitMQ) Pulse() {
	
}

func (this *ConsumerHandlerRabbitMQ)registerMessageProcessor(client cap.IConsumerClient) {
	var onReceive cap.ReceiveHanlder = func (ctx cap.MessageContext) {
		err := this.storeMessage(&ctx)
		if err == nil {
			client.Commit(ctx)
		}
	}

	var onError cap.ErrorHanlder = func (content string) { }

	client.SetOnReceive(onReceive)
	client.SetOnError(onError)
}

func (this *ConsumerHandlerRabbitMQ)storeMessage(context *cap.MessageContext) error {
	if this.ConnectionFactory == nil || this.ConnectionFactory.CreateStorageConnection == nil {
		panic("ConnectionFactory and delegate can not be nil")
	}

	connection, err := this.ConnectionFactory.CreateStorageConnection(this.CapOptions)

	if err != nil {
		return err
	}

	var message = cap.NewCapReceivedMessage(*context)
	message.StatusName = "Scheduled"

	return connection.StoreReceivedMessage(message)
}
