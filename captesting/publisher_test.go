package captesting

import (
	"testing"

	"database/sql"

	"github.com/FeiniuBus/capgo"
	cmysql "github.com/FeiniuBus/capgo/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func Test_Publisher(t *testing.T) {
	options := cap.CapOptions{}
	connectionString := "root:kge2001@tcp(192.168.206.129:3306)/FeiniuCAP?charset=utf8"
	options.UseMySql(connectionString)
	connectionFactory := cap.NewPublisherFactory(CreatePublisher)
	publisher, err := connectionFactory.CreatePublisher(options)
	if err != nil {
		panic(err)
	}
	conn, err := sql.Open("mysql", connectionString)
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

func CreatePublisher(options cap.CapOptions) (cap.IPublisher, error) {
	result := cmysql.NewPublisher(options)
	return result, nil
}
