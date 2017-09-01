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
func NewLockedMessage(message interface{}, messageType int32, dbConnection *sql.DB, dbTransaction *sql.Tx) cap.ILockedMessage {
	lockedMessage := &LockedMessage{
		message:       message,
		messageType:   messageType,
		dbConnection:  dbConnection,
		dbTransaction: dbTransaction,
	}
	lockedMessage.logger = cap.GetLoggerFactory().CreateLogger(lockedMessage)
	return lockedMessage
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
	statement, err := message.dbTransaction.Prepare(query)
	if err != nil {
		message.Rollback()
		message.Dispose()
		message.logger.Log(cap.LevelError, "[Prepare]"+err.Error())
		return nil, err
	}
	return statement, nil
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

func (message *LockedMessage) getMessageId() int {
	if message.messageType == 0 {
		return message.message.(cap.CapPublishedMessage).Id
	} else if message.messageType == 1 {
		return message.message.(cap.CapReceivedMessage).Id
	} else {
		return 0
	}
}

func (message *LockedMessage) logError(err string) {
	message.logger.LogData(cap.LevelError, "[Enqueue]"+err,
		struct {
			MessageType int32
			Message     interface{}
		}{MessageType: message.messageType, Message: message.message})
}

// Enqueue ...
func (message *LockedMessage) Enqueue() (AffectedRows int64, err error) {
	statement := "INSERT INTO `cap.queue` (`MessageId`, `MessageType`) VALUES (?,?);"
	messageId := message.getMessageId()

	if messageId == 0 {
		message.Rollback()
		message.Dispose()
		err := cap.NewCapError("MessageId could not be zero.")
		message.logError(err.Error())
		return 0, err
	}

	result, err := message.dbTransaction.Exec(statement, messageId, message.messageType)
	if err != nil {
		message.Rollback()
		message.Dispose()
		message.logError(err.Error())
		return 0, err
	}

	affectRows, err := result.RowsAffected()
	if err != nil {
		message.Rollback()
		message.Dispose()
		message.logError(err.Error())
		return 0, err
	}

	return affectRows, nil
}
