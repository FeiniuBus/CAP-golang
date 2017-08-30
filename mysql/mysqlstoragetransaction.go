package mysql

import (
	"context"
	"database/sql"

	"github.com/FeiniuBus/capgo"
	_ "github.com/go-sql-driver/mysql"
)

// MySqlStorageTransaction ...
type MySqlStorageTransaction struct {
	Options       *cap.CapOptions
	DbConnection  *sql.DB
	DbTransaction *sql.Tx
}

// NewStorageTransaction ...
func NewStorageTransaction(options *cap.CapOptions) (cap.IStorageTransaction, error) {
	transaction := &MySqlStorageTransaction{}
	transaction.Options = options
	connectionString, err := transaction.Options.GetConnectionString()
	if err != nil {
		return nil, err
	}
	transaction.DbConnection, err = sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	txOptions := &sql.TxOptions{Isolation: sql.LevelReadCommitted}

	transaction.DbTransaction, err = transaction.DbConnection.BeginTx(context.Background(), txOptions)
	if err != nil {
		transaction.Dispose()
		return nil, err
	}
	return transaction, nil
}

// EnqueuePublishedMessage ...
func (transaction *MySqlStorageTransaction) EnqueuePublishedMessage(message *cap.CapPublishedMessage) error {
	statement := "INSERT INTO `cap.queue` (`MessageId`, `MessageType`) VALUES (?,?);"
	result, err := transaction.DbTransaction.Exec(statement, message.Id, 0)
	if err != nil {
		transaction.Dispose()
		return err
	}
	affectRows, err := result.RowsAffected()
	if err != nil {
		transaction.Dispose()
		return err
	}
	if affectRows == 0 {
		transaction.Dispose()
		return cap.NewCapError("EnqueuePublishedMessage : Database execution should affect 1 row but affected 0 row actually.")
	}
	return nil
}

// EnqueueReceivedMessage ...
func (transaction *MySqlStorageTransaction) EnqueueReceivedMessage(message *cap.CapReceivedMessage) error {
	statement := "INSERT INTO `cap.queue` VALUES (?,?);"
	result, err := transaction.DbTransaction.Exec(statement, message.Id, 1)
	if err != nil {
		transaction.Dispose()
		return err
	}
	affectRows, err := result.RowsAffected()
	if err != nil {
		transaction.Dispose()
		return err
	}
	if affectRows == 0 {
		transaction.Dispose()
		return cap.NewCapError("EnqueueReceivedMessage : Database execution should affect 1 row but affected 0 row actually.")
	}
	return nil
}

// UpdatePublishedMessage ...
func (transaction *MySqlStorageTransaction) UpdatePublishedMessage(message *cap.CapPublishedMessage) error {
	statement := "UPDATE `cap.published` "
	statement += " SET `Retries` = ?,`ExpiresAt` = FROM_UNIXTIME(?),`StatusName`= ?"
	statement += " WHERE `Id` = ?"
	result, err := transaction.DbTransaction.Exec(statement, message.Retries, message.ExpiresAt, message.StatusName, message.Id)

	if err != nil {
		transaction.Dispose()
		return err
	}
	affectRows, err := result.RowsAffected()
	if err != nil {
		transaction.Dispose()
		return err
	}
	if affectRows == 0 {
		transaction.Dispose()
		return cap.NewCapError("UpdatePublishedMessage : Database execution should affect 1 row but affected 0 row actually.")
	}
	return nil
}

// UpdateReceivedMessage ...
func (transaction *MySqlStorageTransaction) UpdateReceivedMessage(message *cap.CapReceivedMessage) error {
	statement := "UPDATE `cap.received`"
	statement += " SET `Retries` = ?,`ExpiresAt` = FROM_UNIXTIME(?),`StatusName`= ?"
	statement += " WHERE `Id` = ?"
	result, err := transaction.DbTransaction.Exec(statement, message.Retries, message.ExpiresAt, message.StatusName, message.Id)
	if err != nil {
		transaction.Dispose()
		return err
	}
	affectRows, err := result.RowsAffected()
	if err != nil {
		transaction.Dispose()
		return err
	}
	if affectRows == 0 {
		transaction.Dispose()
		return cap.NewCapError("UpdateReceivedMessage : Database execution should affect 1 row but affected 0 row actually.")
	}
	return nil
}

// Commit ...
func (transaction *MySqlStorageTransaction) Commit() error {
	err := transaction.DbTransaction.Commit()
	if err != nil {
		_ = transaction.DbTransaction.Rollback()
	}
	transaction.Dispose()
	if err != nil {
		return err
	}
	return nil
}

// Dispose ...
func (transaction *MySqlStorageTransaction) Dispose() {
	transaction.DbConnection.Close()
}
