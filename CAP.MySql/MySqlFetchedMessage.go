package cap_mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlFetchedMessage struct{
	messageId int
	messageType int
	dbConnection *sql.DB
	dbTransaction *sql.Tx
}

func NewFetchedMessage(_messageId int, _messageType int, _dbConnection *sql.DB, _dbTransaction *sql.Tx)(*MySqlFetchedMessage){
	result := &MySqlFetchedMessage{}
	result.messageId = _messageId
	result.messageType = _messageType
	result.dbConnection = _dbConnection
	result.dbTransaction = _dbTransaction
	return result
}

func (fetchedMessage *MySqlFetchedMessage) GetMessageId()(messageId int){
	return fetchedMessage.messageId
}

func (fetchedMessage *MySqlFetchedMessage) GetMessageType()(messageType int){
	return fetchedMessage.messageType
}

func (fectchedMessage *MySqlFetchedMessage) RemoveFromQueue() error{
	return nil
}

func (fetchMessage *MySqlFetchedMessage) Requeue() error{
	return nil
}

func (fetchedMessage *MySqlFetchedMessage) Dispose(){
	
}