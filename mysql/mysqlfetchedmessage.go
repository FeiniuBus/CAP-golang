package mysql

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySqlFetchedMessage ...
type MySqlFetchedMessage struct {
	messageId     int
	messageType   int
	dbConnection  *sql.DB
	dbTransaction *sql.Tx
	mutext        *sync.Mutex
	ticker        *time.Ticker
}

// NewFetchedMessage ...
func NewFetchedMessage(_messageId int, _messageType int, _dbConnection *sql.DB, _dbTransaction *sql.Tx) *MySqlFetchedMessage {
	result := &MySqlFetchedMessage{}
	result.messageId = _messageId
	result.messageType = _messageType
	result.dbConnection = _dbConnection
	result.dbTransaction = _dbTransaction
	result.mutext = &sync.Mutex{}
	result.ticker = time.NewTicker(1 * time.Minute)
	go result.keepAlive()
	return result
}

// GetMessageId ...
func (fetchedMessage *MySqlFetchedMessage) GetMessageId() (messageId int) {
	return fetchedMessage.messageId
}

// GetMessageType ...
func (fetchedMessage *MySqlFetchedMessage) GetMessageType() (messageType int) {
	return fetchedMessage.messageType
}

// RemoveFromQueue ...
func (fetchedMessage *MySqlFetchedMessage) RemoveFromQueue() error {
	fetchedMessage.mutext.Lock()
	err := fetchedMessage.dbTransaction.Commit()
	fetchedMessage.mutext.Unlock()
	return err
}

// Requeue ...
func (fetchedMessage *MySqlFetchedMessage) Requeue() error {
	fetchedMessage.mutext.Lock()
	err := fetchedMessage.dbTransaction.Rollback()
	fetchedMessage.mutext.Unlock()
	return err
}

// Dispose ...
func (fetchedMessage *MySqlFetchedMessage) Dispose() error {
	fetchedMessage.mutext.Lock()
	fetchedMessage.ticker.Stop()
	err := fetchedMessage.dbConnection.Close()
	fetchedMessage.mutext.Unlock()

	return err
}

func (fetchedMessage *MySqlFetchedMessage) keepAlive() {
	statement := "SELECT 1;"
	for _ = range fetchedMessage.ticker.C {
		fetchedMessage.mutext.Lock()
		_, _ = fetchedMessage.dbTransaction.Exec(statement)
		fetchedMessage.mutext.Unlock()
	}
}
