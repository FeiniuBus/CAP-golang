package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/FeiniuBus/capgo"
	_ "github.com/go-sql-driver/mysql"
)

// MySqlStorageConnection ...
type MySqlStorageConnection struct {
	Options *cap.CapOptions
	logger  cap.ILogger
}

// NewStorageConnection ...
func NewStorageConnection(options *cap.CapOptions) cap.IStorageConnection {
	connection := &MySqlStorageConnection{}
	connection.Options = options
	connection.logger = cap.GetLoggerFactory().CreateLogger(connection)
	return connection
}

// OpenDbConnection ...
func (connection MySqlStorageConnection) OpenDbConnection() (*sql.DB, error) {
	connectionString, err := connection.Options.GetConnectionString()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[OpenDbConnection]"+err.Error())
		return nil, err
	}
	conn, err := sql.Open("mysql", connectionString)

	if err != nil {
		connection.logger.Log(cap.LevelError, "[OpenDbConnection]"+err.Error())
		return nil, err
	}
	return conn, nil
}

// BeginTransaction ...
func (connection MySqlStorageConnection) BeginTransaction(dbConnection *sql.DB) (*sql.Tx, error) {
	options := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	transaction, err := dbConnection.BeginTx(context.Background(), options)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[BeginTransaction]"+err.Error())
		return nil, err
	}
	return transaction, nil
}

// CreateTransaction ...
func (connection *MySqlStorageConnection) CreateTransaction() (cap.IStorageTransaction, error) {
	transaction, err := NewStorageTransaction(connection.Options)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[CreateTransaction]"+err.Error())
		return nil, err
	}
	return transaction, nil
}

// FetchNextMessage ...
func (connection *MySqlStorageConnection) FetchNextMessage() (cap.IFetchedMessage, error) {
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[FetchNextMessage]"+err.Error())
		return nil, err
	}

	transaction, err := connection.BeginTransaction(conn)
	if err != nil {
		conn.Close()
		connection.logger.Log(cap.LevelError, "[FetchNextMessage]"+err.Error())
		return nil, err
	}

	statement := "SELECT `MessageId`,`MessageType` FROM `cap.queue` LIMIT 1 FOR UPDATE;DELETE FROM `cap.queue` LIMIT 1;"

	row, err := transaction.Query(statement)
	defer row.Close()
	if err != nil {
		conn.Close()
		connection.logger.Log(cap.LevelError, "[FetchNextMessage]"+err.Error())
		return nil, err
	}

	var messageID int
	var messageType int

	if row.Next() == true {
		row.Scan(&messageID, &messageType)
	}

	fetchedMessage := NewFetchedMessage(messageID, messageType, conn, transaction)

	return fetchedMessage, nil
}

// GetFailedPublishedMessages ...
func (connection *MySqlStorageConnection) GetFailedPublishedMessages() ([]*cap.CapPublishedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime,  `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.published` WHERE `StatusName` = 'Failed';"
	conn, err := connection.OpenDbConnection()
	defer conn.Close()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetFailedPublishedMessages]"+err.Error())
		return nil, err
	}

	returnValue := make([]*cap.CapPublishedMessage, 0)

	rows, err := conn.Query(statement)
	defer rows.Close()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetFailedPublishedMessages]"+err.Error())
		return nil, err
	}

	for rows.Next() {
		item := &cap.CapPublishedMessage{}
		err = rows.Scan(&item.Id, &item.Added, &item.Content, &item.ExpiresAt, &item.LastWarnedTime, &item.MessageId, &item.Name, &item.Retries, &item.StatusName, &item.TransactionId)
		if err != nil {
			connection.logger.Log(cap.LevelError, "[GetFailedPublishedMessages]"+err.Error())
			return nil, err
		}
		returnValue = append(returnValue, item)
	}
	return returnValue, nil
}

// GetFailedReceivedMessages ...
func (connection *MySqlStorageConnection) GetFailedReceivedMessages() ([]*cap.CapReceivedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, `Group`, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.received` WHERE `StatusName` = 'Failed';"
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetFailedPublishedMessages]"+err.Error())
		return nil, err
	}
	defer conn.Close()
	returnValue := make([]*cap.CapReceivedMessage, 0)

	rows, err := conn.Query(statement)

	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetFailedPublishedMessages]"+err.Error())
		return nil, err
	}

	for rows.Next() {
		item := &cap.CapReceivedMessage{}
		err = rows.Scan(&item.Id, &item.Added, &item.Content, &item.ExpiresAt, &item.Group, &item.LastWarnedTime, &item.MessageId, &item.Name, &item.Retries, &item.StatusName, &item.TransactionId)
		if err != nil {
			connection.logger.Log(cap.LevelError, "[GetFailedPublishedMessages]"+err.Error())
			return nil, err
		}
		returnValue = append(returnValue, item)
	}
	return returnValue, nil
}

// GetNextPublishedMessageToBeEnqueued ...
func (connection *MySqlStorageConnection) GetNextPublishedMessageToBeEnqueued() (*cap.CapPublishedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.published` WHERE `StatusName` = 'Scheduled' LIMIT 1;"
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetNextPublishedMessageToBeEnqueued]"+err.Error())
		return nil, err
	}
	defer conn.Close()
	if conn == nil {
		err = cap.NewCapError("Database connection is nil.")
		connection.logger.Log(cap.LevelError, "[GetNextPublishedMessageToBeEnqueued]"+err.Error())
		return nil, err
	}

	rows, err := conn.Query(statement)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetNextPublishedMessageToBeEnqueued]"+err.Error())
		return nil, err
	}
	defer rows.Close()

	message := &cap.CapPublishedMessage{}
	if rows.Next() {
		rows.Scan(&message.Id, &message.Added, &message.Content, &message.ExpiresAt, &message.LastWarnedTime, &message.MessageId, &message.Name, &message.Retries, &message.StatusName, &message.TransactionId)
	}

	return message, nil
}

// GetNextReceviedMessageToBeEnqueued ..
func (connection *MySqlStorageConnection) GetNextReceviedMessageToBeEnqueued() (*cap.CapReceivedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, `Group`, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.received` WHERE `StatusName` = 'Scheduled' LIMIT 1;"
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetNextReceviedMessageToBeEnqueued]"+err.Error())
		return nil, err
	}
	defer conn.Close()
	rows, err := conn.Query(statement)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetNextReceviedMessageToBeEnqueued]"+err.Error())
		return nil, err
	}
	defer rows.Close()
	message := &cap.CapReceivedMessage{}

	if rows.Next() {
		rows.Scan(&message.Id, &message.Added, &message.Content, &message.ExpiresAt, &message.Group, &message.LastWarnedTime, &message.MessageId, &message.Name, &message.Retries, &message.StatusName, &message.TransactionId)
	}

	return message, nil
}

// GetPublishedMessage ...
func (connection *MySqlStorageConnection) GetPublishedMessage(id int) (*cap.CapPublishedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.published` WHERE `Id`=?;"
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetPublishedMessage]"+err.Error())
		return nil, err
	}
	defer conn.Close()
	rows, err := conn.Query(statement, id)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetPublishedMessage]"+err.Error())
		return nil, err
	}
	defer rows.Close()
	message := &cap.CapPublishedMessage{}

	if rows.Next() {
		rows.Scan(&message.Id, &message.Added, &message.Content, &message.ExpiresAt, &message.LastWarnedTime, &message.MessageId, &message.Name, &message.Retries, &message.StatusName, &message.TransactionId)
	}

	return message, nil
}

// GetReceivedMessage ...
func (connection *MySqlStorageConnection) GetReceivedMessage(id int) (*cap.CapReceivedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, `Group`, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.received` WHERE `Id`=?;"
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetReceivedMessage]"+err.Error())
		return nil, err
	}
	defer conn.Close()
	rows, err := conn.Query(statement, id)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetReceivedMessage]"+err.Error())
		return nil, err
	}
	defer rows.Close()
	message := &cap.CapReceivedMessage{}

	if rows.Next() {
		rows.Scan(&message.Id, &message.Added, &message.Content, &message.ExpiresAt, &message.Group, &message.LastWarnedTime, &message.MessageId, &message.Name, &message.Retries, &message.StatusName, &message.TransactionId)
	}

	return message, nil
}

// GetNextLockedMessageToBeEnqueued ...
func (connection *MySqlStorageConnection) GetNextLockedMessageToBeEnqueued(messageType int32) (cap.ILockedMessage, error) {
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.LogData(cap.LevelError, "[GetNextLockedMessageToBeEnqueued]"+err.Error(), struct{ MessageType int32 }{MessageType: messageType})
		return nil, err
	}

	transaction, err := connection.BeginTransaction(conn)
	if err != nil {
		conn.Close()
		connection.logger.LogData(cap.LevelError, "[GetNextLockedMessageToBeEnqueued]"+err.Error(), struct{ MessageType int32 }{MessageType: messageType})
		return nil, err
	}
	var message cap.ILockedMessage
	if messageType == 0 {
		message, err = connection.getNextPublishedLockedMessageToBeEnqueued(conn, transaction)
		if err != nil {
			transaction.Rollback()
			conn.Close()
			connection.logger.LogData(cap.LevelError, "[GetNextLockedMessageToBeEnqueued]"+err.Error(), struct{ MessageType int32 }{MessageType: messageType})
			return nil, err
		}
	} else if messageType == 1 {
		message, err = connection.getNextReceivedLockedMessageToBeEnqueued(conn, transaction)
		if err != nil {
			transaction.Rollback()
			conn.Close()
			connection.logger.LogData(cap.LevelError, "[GetNextLockedMessageToBeEnqueued]"+err.Error(), struct{ MessageType int32 }{MessageType: messageType})
			return nil, err
		}
	} else {
		err = cap.NewCapError("Unknown MessageType.")
		transaction.Rollback()
		conn.Close()
		connection.logger.LogData(cap.LevelError, "[GetNextLockedMessageToBeEnqueued]"+err.Error(), struct{ MessageType int32 }{MessageType: messageType})
		return nil, err
	}

	if message == nil {
		transaction.Rollback()
		conn.Close()
		return nil, nil
	}

	return message, nil
}

func (connection *MySqlStorageConnection) getNextReceivedLockedMessageToBeEnqueued(conn *sql.DB, transaction *sql.Tx) (cap.ILockedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, `Group`, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.received` WHERE `StatusName` = 'Scheduled' LIMIT 1 FOR UPDATE;"

	rows, err := transaction.Query(statement)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[getNextReceivedLockedMessageToBeEnqueued]"+err.Error())
		return nil, err
	}
	defer rows.Close()
	message := &cap.CapReceivedMessage{}

	if rows.Next() {
		rows.Scan(&message.Id, &message.Added, &message.Content, &message.ExpiresAt, &message.Group, &message.LastWarnedTime, &message.MessageId, &message.Name, &message.Retries, &message.StatusName, &message.TransactionId)
	} else {
		return nil, nil
	}

	if message.Id == 0 {
		return nil, nil
	}

	return NewLockedMessage(message, 1, conn, transaction, connection.Options), nil
}

func (connection *MySqlStorageConnection) getNextPublishedLockedMessageToBeEnqueued(conn *sql.DB, transaction *sql.Tx) (cap.ILockedMessage, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.published` WHERE `StatusName` = 'Scheduled' LIMIT 1 FOR UPDATE;"

	rows, err := transaction.Query(statement)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[getNextPublishedLockedMessageToBeEnqueued]"+err.Error())
		return nil, err
	}
	defer rows.Close()
	message := &cap.CapPublishedMessage{}

	if rows.Next() {
		rows.Scan(&message.Id, &message.Added, &message.Content, &message.ExpiresAt, &message.LastWarnedTime, &message.MessageId, &message.Name, &message.Retries, &message.StatusName, &message.TransactionId)
	} else {
		return nil, nil
	}

	if message.Id == 0 {
		return nil, nil
	}

	return NewLockedMessage(message, 0, conn, transaction, connection.Options), nil
}

// StoreReceivedMessage ...
func (connection *MySqlStorageConnection) StoreReceivedMessage(message *cap.CapReceivedMessage) error {
	statement := "INSERT INTO `cap.received`(`Name`,`Group`,`Content`,`Retries`,`Added`,`ExpiresAt`,`StatusName`,`MessageId`,`TransactionId`)"
	statement += " VALUES(?,?,?,?,?,?,?,?,?);"
	conn, err := connection.OpenDbConnection()
	defer conn.Close()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[StoreReceivedMessage]"+err.Error())
		return err
	}

	feiniuMessage := cap.FeiniuBusMessage{
		MetaData: cap.FeiniuBusMessageMetaData{},
	}
	err = json.Unmarshal([]byte(message.Content), &feiniuMessage)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[StoreReceivedMessage]"+err.Error())
		return err
	}

	result, err := conn.Exec(statement, message.Name, message.Group, message.Content, message.Retries, time.Now(), nil, message.StatusName, feiniuMessage.MetaData.MessageID, feiniuMessage.MetaData.TransactionID)
	if err != nil {
		connection.logger.Log(cap.LevelError, "[StoreReceivedMessage]"+err.Error())
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[StoreReceivedMessage]"+err.Error())
		return err
	}
	if rowAffected == int64(0) {
		err = cap.NewCapError("StoreReceivedMessage : Database execution should affect 1 row but affected 0 row actually.")
		connection.logger.Log(cap.LevelError, "[StoreReceivedMessage]"+err.Error())
		return err
	}
	return nil
}

// GetFailedPublishedLockedMessages ...
func (connection *MySqlStorageConnection) GetFailedPublishedLockedMessages(conn *sql.DB, transaction *sql.Tx) (cap.ILockedMessages, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime,  `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.published` WHERE `StatusName` = 'Failed' FOR UPDATE;"

	returnValue := NewLockedMessages(0, conn, transaction, connection.Options)

	rows, err := transaction.Query(statement)
	defer rows.Close()
	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetFailedPublishedLockedMessages]"+err.Error())
		return nil, err
	}

	for rows.Next() {
		item := &cap.CapPublishedMessage{}
		err = rows.Scan(&item.Id, &item.Added, &item.Content, &item.ExpiresAt, &item.LastWarnedTime, &item.MessageId, &item.Name, &item.Retries, &item.StatusName, &item.TransactionId)
		if err != nil {
			connection.logger.Log(cap.LevelError, "[GetFailedPublishedLockedMessages]"+err.Error())
			return nil, err
		}
		returnValue.AppendMessage(item)
	}

	return returnValue, nil
}

// GetFailedReceivedMessages ...
func (connection *MySqlStorageConnection) GetFailedReceivedLockedMessages(conn *sql.DB, transaction *sql.Tx) (cap.ILockedMessages, error) {
	statement := "SELECT `Id`, CONVERT(UNIX_TIMESTAMP(`Added`),SIGNED) AS Added, `Content`, CONVERT(UNIX_TIMESTAMP(`ExpiresAt`),SIGNED) AS ExpiresAt, `Group`, CONVERT(UNIX_TIMESTAMP(`LastWarnedTime`),SIGNED) AS LastWarnedTime, `MessageId`, `Name`, `Retries`, `StatusName`, `TransactionId` FROM `cap.received` WHERE `StatusName` = 'Failed' FOR UPDATE;"
	returnValue := NewLockedMessages(1, conn, transaction, connection.Options)

	rows, err := conn.Query(statement)

	if err != nil {
		connection.logger.Log(cap.LevelError, "[GetFailedReceivedLockedMessages]"+err.Error())
		return nil, err
	}

	for rows.Next() {
		item := &cap.CapReceivedMessage{}
		err = rows.Scan(&item.Id, &item.Added, &item.Content, &item.ExpiresAt, &item.Group, &item.LastWarnedTime, &item.MessageId, &item.Name, &item.Retries, &item.StatusName, &item.TransactionId)
		if err != nil {
			connection.logger.Log(cap.LevelError, "[GetFailedReceivedLockedMessages]"+err.Error())
			return nil, err
		}
		returnValue.AppendMessage(item)
	}
	return returnValue, nil
}

// GetFailedLockedMessages ...
func (connection *MySqlStorageConnection) GetFailedLockedMessages(messageType int32) (cap.ILockedMessages, error) {
	conn, err := connection.OpenDbConnection()
	if err != nil {
		connection.logger.LogData(cap.LevelError, err.Error(), struct{ MessageType int32 }{MessageType: messageType})
		return nil, err
	}
	transaction, err := connection.BeginTransaction(conn)
	if err != nil {
		conn.Close()
		connection.logger.LogData(cap.LevelError, err.Error(), struct{ MessageType int32 }{MessageType: messageType})
		return nil, err
	}
	if messageType == 0 {
		messages, err := connection.GetFailedPublishedLockedMessages(conn, transaction)
		if err != nil {
			transaction.Rollback()
			conn.Close()
			connection.logger.LogData(cap.LevelError, err.Error(), struct{ MessageType int32 }{MessageType: messageType})
			return nil, err
		}
		return messages, nil
	} else if messageType == 1 {
		messages, err := connection.GetFailedReceivedLockedMessages(conn, transaction)
		if err != nil {
			transaction.Rollback()
			conn.Close()
			connection.logger.LogData(cap.LevelError, err.Error(), struct{ MessageType int32 }{MessageType: messageType})
			return nil, err
		}
		return messages, nil
	} else {
		err = cap.NewCapError("Unknown MessageType.")
		transaction.Rollback()
		conn.Close()
		connection.logger.LogData(cap.LevelError, err.Error(), struct{ MessageType int32 }{MessageType: messageType})
		return nil, err
	}
}
