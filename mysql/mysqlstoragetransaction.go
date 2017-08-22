package mysql

import(
	"github.com/FeiniuBus/capgo"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlStorageTransaction struct
{
	Options cap.CapOptions
	DbConnection *sql.DB;
	DbTransaction *sql.Tx;
}

func NewStorageTransaction(options cap.CapOptions) (cap.IStorageTransaction, error) {
	transaction := &MySqlStorageTransaction{}
	transaction.Options = options
	connectionString, err := transaction.Options.GetConnectionString()
	if err != nil{
		return nil, err
	}
	transaction.DbConnection, err = sql.Open("mysql",connectionString)
	if err != nil{
		return nil, err
	}
	transaction.DbTransaction, err = transaction.DbConnection.Begin()
	if err != nil{
		return nil, err
	}
	return transaction,nil
}

func (transaction *MySqlStorageTransaction) EnqueuePublishedMessage(message *cap.CapPublishedMessage) error{
	statement := "INSERT INTO "+ cap.QueueTableName +" VALUES (?,?)"
	stmt,err := transaction.DbConnection.Prepare(statement)
	if err != nil{
		return err
	}
	_, err=stmt.Exec(message.Id,0)
	if err != nil{
		return err
	} 
	return nil
}

func (transaction *MySqlStorageTransaction) EnqueueReceivedMessage(message *cap.CapReceivedMessage) error{
	statement := "INSERT INTO cap.queue VALUES (?,?)"
	stmt,err := transaction.DbConnection.Prepare(statement)
	if err != nil{
		return err
	}
	_, err=stmt.Exec(message.Id,1)
	if err != nil{
		return err
	} 
	return nil 
}

func (transaction *MySqlStorageTransaction) UpdatePublishedMessage(message *cap.CapPublishedMessage) error{
	statement := "UPDATE " + cap.PublishedTableName 
	statement += " SET `Retries` = ?,`ExpiresAt` = ?,`StatusName`= ?"
	statement += " WHERE `Id` = ?"
	stmt,err := transaction.DbConnection.Prepare(statement)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(message.Retries, message.ExpiresAt, message.Id)
	if err != nil{
		return err
	}
	return nil
}

func (transaction *MySqlStorageTransaction) UpdateReceivedMessage(message *cap.CapReceivedMessage) error{
	statement := "UPDATE " + cap.ReceivedTableName 
	statement += " SET `Retries` = ?,`ExpiresAt` = ?,`StatusName`= ?"
	statement += " WHERE `Id` = ?"
	stmt,err := transaction.DbConnection.Prepare(statement)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(message.Retries, message.ExpiresAt, message.Id)
	if err != nil{
		return err
	}
	return nil
}

func (transaction *MySqlStorageTransaction) Commit() error{
	err := transaction.DbTransaction.Commit()
	if err != nil{
		return err
	}
	return nil
}