package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type MySqlFetchedMessage struct{
	messageId int
	messageType int
	dbConnection *sql.DB
	dbTransaction *sql.Tx
	mutext *sync.Mutex
}

func NewFetchedMessage(_messageId int, _messageType int, _dbConnection *sql.DB, _dbTransaction *sql.Tx)(*MySqlFetchedMessage){
	result := &MySqlFetchedMessage{}
	result.messageId = _messageId
	result.messageType = _messageType
	result.dbConnection = _dbConnection
	result.dbTransaction = _dbTransaction
	result.mutext = &sync.Mutex{}
	go result.keeyAlive()
	return result
}

func (fetchedMessage *MySqlFetchedMessage) GetMessageId()(messageId int){
	return fetchedMessage.messageId
}

func (fetchedMessage *MySqlFetchedMessage) GetMessageType()(messageType int){
	return fetchedMessage.messageType
}

func (fetchedMessage *MySqlFetchedMessage) RemoveFromQueue() error{
	fetchedMessage.mutext.Lock()
	err := fetchedMessage.dbTransaction.Commit()
	fetchedMessage.mutext.Unlock()
	return err
}

func (fetchedMessage *MySqlFetchedMessage) Requeue() error{
	fetchedMessage.mutext.Lock()
	err := fetchedMessage.dbTransaction.Rollback()
	fetchedMessage.mutext.Unlock()
	return err
}

func (fetchedMessage *MySqlFetchedMessage) Dispose() error{
	fetchedMessage.mutext.Lock()
	err := fetchedMessage.dbConnection.Close()
	fetchedMessage.mutext.Unlock()
	return err
}

func (fetchedMessage *MySqlFetchedMessage) keeyAlive(){
	statement := "SELECT 1"
	tick := time.Tick(60 * time.Second)
	for{
		select{
		case <- tick :
			fetchedMessage.mutext.Lock()
			_, _ = fetchedMessage.dbTransaction.Exec(statement)
			fetchedMessage.mutext.Unlock()
		default:
			
		}
	}
}