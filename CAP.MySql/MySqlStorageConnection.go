package cap_mysql

import(
	"../CAP"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlStorageConnection struct{

}

func (connection *MySqlStorageConnection) CreateTransaction(dbConnection interface{})(cap.IStorageTransaction,error){
	transaction := new(MySqlStorageTransaction)
	transaction.DbConnection = dbConnection.(*sql.DB)
	return transaction,nil
}

func (connection *MySqlStorageConnection) FetchNextMessage() (cap.IFetchedMessage,error){
	conn, err := sql.Open("mysql","CapConnectionString")
	if err != nil{
		return nil,err
	}

	transaction,err := conn.Begin()
	if err != nil{
		return nil,err
	}

	statement := "SELECT `MessageId`,`MessageType` FROM `{_prefix}.queue` LIMIT 1 FOR UPDATE;"
	statement += "DELETE FROM `{_prefix}.queue` LIMIT 1;"

	row, err := transaction.Query(statement)
	if err != nil{
		return nil, err
	}

	var messageId int
	var messageType int

	if row.Next() == true {
		row.Scan(&messageId, &messageType)
	}else{
		return nil,nil
	}
	
	fetchedMessage := &MySqlFetchedMessage{}
	fetchedMessage.dbConnection = conn
	fetchedMessage.dbTransaction = transaction


	return fetchedMessage,nil 
}