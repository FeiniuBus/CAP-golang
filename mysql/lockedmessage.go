package mysql

import (
	"database/sql"

	cap "github.com/FeiniuBus/capgo"
)

// LockedMessage ...
type LockedMessage struct {
	message       interface{}
	messageType   int32
	dbConnection  *sql.DB
	dbTransaction *sql.Tx
	logger        cap.ILogger
}

// NewLockedMessage ...
func NewLockedMessage(message interface{}, messageType int32, dbConnection *sql.DB, dbTransaction *sql.Tx) ILockedMessage {
	message := &LockedMessage{
		message:       message,
		messageType:   messageType,
		dbConnection:  dbConnection,
		dbTransaction: dbTransaction,
	}
	message.logger = cap.GetLoggerFactory().CreateLogger(message)
	return message
}

// GetMessage ...
func (message *LockedMessage) GetMessage() interface{} {
	return message.message
}

// GetMessageType ...
func (message *LockedMessage) GetMessageType() int32 {
	return message.messageType
}

// Prepare ...
func (message *LockedMessage) Prepare(query string) (stmt interface{}, err error) {
	stmt, err := message.dbTransaction.Prepare(query)
	if err != nil {
		message.Rollback()
		message.Dispose()
		message.logger.Log(cap.LevelError, "[Prepare]"+err.Error())
		return err
	}
	return stmt, nil
}

// Commit ...
func (message *LockedMessage) Commit() error {
	err := message.dbTransaction.Commit()
	if err != nil {
		message.logger.Log(cap.LevelError, err.Error())
		return err
	}
	return nil
}

// Rollback ...
func (message *LockedMessage) Rollback() error {
	err := message.dbTransaction.Rollback()
	if err != nil {
		message.logger.Log(cap.LevelError, err.Error())
		return err
	}
	return nil
}

// Dispose ...
func (message *LockedMessage) Dispose() {
	err := message.dbConnection.Close()
	if err != nil {
		message.logger.Log(cap.LevelError, "[Dispose]"+err.Error())
	}
}
