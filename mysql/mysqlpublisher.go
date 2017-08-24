package mysql

import (
	"database/sql"
	"time"

	"github.com/FeiniuBus/capgo"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlPublisher struct {
	Options          *cap.CapOptions
	IsCapOpenedTrans bool
}

func NewPublisher(options *cap.CapOptions) cap.IPublisher {
	pubisher := &MySqlPublisher{}
	pubisher.Options = options
	return pubisher
}

func (publihser *MySqlPublisher) Publish(name string, content string, connection interface{}, transaction interface{}) error {
	var dbConnection *sql.DB
	var dbTransaction *sql.Tx

	if connection == nil {
		connectionString, err := publihser.Options.GetConnectionString()
		if err != nil {
			return err
		}
		dbConnection, err = sql.Open("mysql", connectionString)
		if err != nil {
			return err
		}
		dbTransaction, err = dbConnection.Begin()
		if err != nil {
			return err
		}
		publihser.IsCapOpenedTrans = true
	} else {
		dbConnection = connection.(*sql.DB)
		dbTransaction = transaction.(*sql.Tx)
	}
	statement := "INSERT INTO `cap.published` (`Name`,`Content`,`Retries`,`Added`,`ExpiresAt`,`StatusName`)"
	statement += "VALUES(?,?,?,?,?,?)"
	result, err := dbTransaction.Exec(statement, name, content, 0, time.Now(), nil, "Scheduled")
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == int64(0) {
		return cap.NewCapError("Database execution should affect 1 row but affected 0 row actually.")
	}

	if publihser.IsCapOpenedTrans {
		err = dbTransaction.Commit()
		err = dbConnection.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
