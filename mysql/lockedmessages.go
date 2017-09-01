package mysql

import (
	"database/sql"

	cap "github.com/FeiniuBus/capgo"
)

// LockedMessages ...
type LockedMessages struct {
	messages      []cap.ILockedMessage
	messageType   int32
	dbConnection  *sql.DB
	dbTransaction *sql.Tx
	logger        cap.ILogger
	capOptions    *cap.CapOptions
}

// NewLockedMessages ...
func NewLockedMessages(messages []interface{}, messageType int32, dbConnection *sql.DB, dbTransaction *sql.Tx, capOptions *cap.CapOptions) cap.ILockedMessages {

	if messages == nil || len(messages) == 0 {
		return nil
	}

	lockedMessages := &LockedMessages{
		messageType:   messageType,
		dbConnection:  dbConnection,
		dbTransaction: dbTransaction,
		capOptions:    capOptions,
	}
	lockedMessages.logger = cap.GetLoggerFactory().CreateLogger(lockedMessages)

	for _, val := range messages {
		lockedMessages.messages = append(lockedMessages.messages,
			NewLockedMessage(val,
				lockedMessages.messageType,
				lockedMessages.dbConnection,
				lockedMessages.dbTransaction,
				lockedMessages.capOptions))

	}

	return lockedMessages
}

// GetMessages ...
func (message *LockedMessages) GetMessages() []cap.ILockedMessage {
	return message.messages
}

// GetMessageType ...
func (message *LockedMessages) GetMessageType() int32 {
	return message.messageType
}

// Prepare ...
func (message *LockedMessages) Prepare(query string) (stmt interface{}, err error) {
	statement, err := message.dbTransaction.Prepare(query)
	if err != nil {
		message.logger.Log(cap.LevelError, "[Prepare]"+err.Error())
		return nil, err
	}
	return statement, nil
}

// Commit ...
func (message *LockedMessages) Commit() error {
	err := message.dbTransaction.Commit()
	if err != nil {
		message.logger.Log(cap.LevelError, err.Error())
		return err
	}
	return nil
}

// Rollback ...
func (message *LockedMessages) Rollback() error {
	err := message.dbTransaction.Rollback()
	if err != nil {
		message.logger.Log(cap.LevelError, err.Error())
		return err
	}
	return nil
}

// Dispose ...
func (message *LockedMessages) Dispose() {
	err := message.dbConnection.Close()
	if err != nil {
		message.logger.Log(cap.LevelError, "[Dispose]"+err.Error())
	}
}

func (message *LockedMessages) logError(err string) {
	message.logger.LogData(cap.LevelError, "[Enqueue]"+err,
		struct {
			MessageType int32
			Messages    []cap.ILockedMessage
		}{MessageType: message.messageType, Messages: message.messages})
}

// ChangeStates ...
func (message *LockedMessages) ChangeStates(state cap.IState) error {
	for _, val := range message.GetMessages() {
		err := val.ChangeState(state)
		if err != nil {
			message.logError(err.Error())
			return err
		}
	}
	return nil
}
