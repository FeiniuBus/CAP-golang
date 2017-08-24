package main

import (
	"database/sql"
	"github.com/FeiniuBus/capgo"
	cmysql "github.com/FeiniuBus/capgo/mysql"
	_ "github.com/go-sql-driver/mysql" 
)

var CapOptions cap.CapOptions
var ConnectionString string
var PublisherFactory *cap.PublisherFactory

func CreatePublisher(options cap.CapOptions) (cap.IPublisher, error) {
	result := cmysql.NewPublisher(options)
	return result, nil
}

func init(){
	CapOptions = cap.CapOptions{}
	ConnectionString = "root:kge2001@tcp(192.168.206.129:3306)/FeiniuCAP?charset=utf8"
	CapOptions.UseMySql(ConnectionString)
	PublisherFactory = cap.NewPublisherFactory(CreatePublisher) 
}

func main(){
	publisher, err := PublisherFactory.CreatePublisher(CapOptions)
	if err != nil {
		panic(err)
	}
	conn, err := sql.Open("mysql", ConnectionString)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	trans, err := conn.Begin()
	if err != nil {
		panic(err)
	}
	err = publisher.Publish("test", "test", conn, trans)
	if err != nil {
		panic(err)
	}
}