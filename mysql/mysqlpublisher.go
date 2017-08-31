package mysql

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/FeiniuBus/capgo"
	_ "github.com/go-sql-driver/mysql"
)

// MySqlPublisher ..
type MySqlPublisher struct {
	Options          *cap.CapOptions
	IsCapOpenedTrans bool
	logger           cap.ILogger
}

// NewPublisher ..
func NewPublisher(options *cap.CapOptions) cap.IPublisher {
	pubisher := &MySqlPublisher{}
	pubisher.Options = options
	pubisher.logger = cap.GetLoggerFactory().CreateLogger(pubisher)
	return pubisher
}

// Publish ..
func (publisher *MySqlPublisher) Publish(descriptors []*cap.MessageDescriptor, connection interface{}, transaction interface{}) error {
	if len(descriptors) == 0 {
		return nil
	}

	var dbConnection *sql.DB
	var dbTransaction *sql.Tx

	if connection == nil {
		connectionString, err := publisher.Options.GetConnectionString()
		if err != nil {
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
		dbConnection, err = sql.Open("mysql", connectionString)
		if err != nil {
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
		dbTransaction, err = dbConnection.Begin()
		if err != nil {
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
		publisher.IsCapOpenedTrans = true
	} else {
		dbConnection = connection.(*sql.DB)
		dbTransaction = transaction.(*sql.Tx)
	}

	statement := "INSERT INTO `cap.published` (`Name`,`Content`,`Retries`,`Added`,`ExpiresAt`,`StatusName`,`MessageId`,`TransactionId`)"
	statement += "VALUES(?,?,?,?,?,?,?,?)"

	transactionID := cap.NewID()

	for _, val := range descriptors {
		jsonStr, err := json.Marshal(val.Content)
		if err != nil {
			if publisher.IsCapOpenedTrans {
				_ = dbTransaction.Rollback()
				_ = dbConnection.Close()
			}
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
		feiniuMessage := cap.FeiniuBusMessage{
			MetaData: cap.FeiniuBusMessageMetaData{
				TransactionID: transactionID,
				MessageID:     cap.NewID(),
			},
			Content: string(jsonStr),
		}

		messageContent, err := json.Marshal(feiniuMessage)
		messageStr := string(messageContent)
		if err != nil {
			if publisher.IsCapOpenedTrans {
				_ = dbTransaction.Rollback()
				_ = dbConnection.Close()
			}
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}

		result, err := dbTransaction.Exec(statement, val.Name, messageStr, 0, time.Now(), nil, "Scheduled", feiniuMessage.MetaData.MessageID, transactionID)

		if err != nil {
			if publisher.IsCapOpenedTrans {
				_ = dbTransaction.Rollback()
				_ = dbConnection.Close()
			}
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
		affectedRows, err := result.RowsAffected()
		if err != nil {
			if publisher.IsCapOpenedTrans {
				_ = dbTransaction.Rollback()
				_ = dbConnection.Close()
			}
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
		if affectedRows == int64(0) {
			if publisher.IsCapOpenedTrans {
				_ = dbTransaction.Rollback()
				_ = dbConnection.Close()
			}
			err = cap.NewCapError("Database execution should affect 1 row but affected 0 row actually.")
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
	}

	if publisher.IsCapOpenedTrans {
		err := dbTransaction.Commit()
		_ = dbConnection.Close()
		if err != nil {
			publisher.logger.Log(cap.LevelError, "[Publish]"+err.Error())
			return err
		}
	}

	return nil
}

// PublishOne ...
func (publisher *MySqlPublisher) PublishOne(name string, content interface{}, connection interface{}, transaction interface{}) error {
	descriptors := make([]*cap.MessageDescriptor, 0)
	descriptors = append(descriptors, &cap.MessageDescriptor{
		Name:    name,
		Content: content,
	})
	err := publisher.Publish(descriptors, connection, transaction)
	if err != nil {
		publisher.logger.Log(cap.LevelError, "[PublishOne]"+err.Error())
	}
	return err
}
