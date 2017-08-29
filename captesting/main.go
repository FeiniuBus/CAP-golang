package main

import (
	"database/sql"

	"github.com/FeiniuBus/capgo"
	cmysql "github.com/FeiniuBus/capgo/mysql"
	crabbitmq "github.com/FeiniuBus/capgo/rabbitmq"
	_ "github.com/go-sql-driver/mysql"
)

var CapOptions *cap.CapOptions
var RabbitMQOptions *crabbitmq.RabbitMQOptions
var ConnectionString string
var PublisherFactory *cap.PublisherFactory
var StorageConnectionFactory *cap.StorageConnectionFactory
var ProcessorServer *cap.ProcessorServer
var CallbackRegister *cap.CallbackRegister
var Bootstrapper *cap.Bootstrapper

func CreatePublisher(options *cap.CapOptions) (cap.IPublisher, error) {
	result := cmysql.NewPublisher(options)
	return result, nil
}

func CreateStorageConnection(options *cap.CapOptions) (cap.IStorageConnection, error) {
	result := cmysql.NewStorageConnection(options)
	return result, nil
}

func init() {
	CapOptions = cap.NewCapOptions()
	ConnectionString = "root:kge2001@tcp(192.168.206.129:3306)/FeiniuCAP?charset=utf8"
	CapOptions.ConnectionString = ConnectionString
	ProcessorServer = cap.NewProcessorServer()
	ProcessorServer.Container.Register(cap.NewFailedJobProcessor(CapOptions, StorageConnectionFactory))
	ProcessorServer.Container.Register(cap.NewPublishQueuer(CapOptions, StorageConnectionFactory))
	ProcessorServer.Container.Register(cap.NewSubscribeQueuer(CapOptions, StorageConnectionFactory))
	ProcessorServer.Container.Register(cap.NewDefaultDispatcher(CapOptions, StorageConnectionFactory))

	Bootstrapper = cap.NewBootstrapper(CapOptions, StorageConnectionFactory)
	Bootstrapper.Servers = append(Bootstrapper.Servers,ProcessorServer)

	RabbitMQOptions = crabbitmq.RabbitMQConfig
	RabbitMQOptions.SetHostName("192.168.206.128")
	RabbitMQOptions.SetUserName("andrew")
	RabbitMQOptions.SetPassword("kge2001")
	crabbitmq.Prepare(Bootstrapper, *RabbitMQOptions)

	PublisherFactory = cap.NewPublisherFactory(CreatePublisher)
	StorageConnectionFactory = cap.NewStorageConnectionFactory(CreateStorageConnection)

	publisher, err := PublisherFactory.CreatePublisher(CapOptions)
	if err != nil {
		panic(err)
	}

	dbConnection, err := sql.Open("mysql", ConnectionString)
	defer dbConnection.Close()
	if err != nil {
		panic(err)
	}

	dbTransaction, err := dbConnection.Begin()
	if err != nil {
		panic(err)
	}

	err = publisher.Publish("test", "test", dbConnection, dbTransaction)
	if err != nil {
		panic(err)
	}
	err = dbTransaction.Commit()
	if err != nil {
		panic(err)
	}
}

func main() {
	go Bootstrapper.Bootstrap()

	for {

	}
	//Bootstrapper.WaitForTerminalSignal()
}
