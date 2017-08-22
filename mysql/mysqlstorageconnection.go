package mysql

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlStorageConnection struct{
	Options cap.CapOptions
}

func NewStorageConnection(options cap.CapOptions) cap.IStorageConnection{
	connection := &MySqlStorageConnection{}
	connection.Options = options
	return connection
}

func (connection MySqlStorageConnection) OpenDbConnection() (*sql.DB, error){
	connectionString, err := connection.Options.GetConnectionString()
	if err != nil{
		return nil, err
	}
	conn, err := sql.Open("mysql",connectionString)
	if err != nil{
		return nil,err
	} 
	return conn,nil
}

func (connection *MySqlStorageConnection) CreateTransaction()(cap.IStorageTransaction,error){
	transaction, err := NewStorageTransaction(connection.Options)
	if err != nil {
		return nil,err
	}
	return transaction, nil
}

func (connection *MySqlStorageConnection) FetchNextMessage() (cap.IFetchedMessage,error){
	conn, err := connection.OpenDbConnection()
	if err != nil{
		return nil,err
	}

	transaction,err := conn.Begin()
	if err != nil{
		return nil,err
	}

	statement := "SELECT `MessageId`,`MessageType` FROM `cap.queue` LIMIT 1 FOR UPDATE;"
	statement += "DELETE FROM `cap.queue` LIMIT 1;"

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

func (connection *MySqlStorageConnection) GetFailedPublishedMessages() ([]*cap.CapPublishedMessage,error){
	statement := "SELECT * FROM `cap.published` WHERE `StatusName` = 'Failed';"
	conn, err := connection.OpenDbConnection()
	if err != nil{
		return nil, err
	}

	returnValue := make([]*cap.CapPublishedMessage,0)

	rows, err := conn.Query(statement)

	if err != nil{
		return nil, err
	}

	for rows.Next() {
		item := &cap.CapPublishedMessage{}
		err = rows.Scan(&item)
		if err != nil{
			return nil, err
		}
		returnValue = append(returnValue, item)
	}
	return returnValue , nil
}

func (connection *MySqlStorageConnection) GetNextPublishedMessageToBeEnqueued() (*cap.CapPublishedMessage,error){
	statement := "SELECT * FROM `cap.published` WHERE `StatusName` = 'Scheduled' LIMIT 1;"
	conn, err := connection.OpenDbConnection()
	if err != nil {
		return nil, err
	}
	rows,err := conn.Query(statement)
	if err!=nil{
		return nil, err
	}
	message := &cap.CapPublishedMessage{}
	if rows.Next(){
		rows.Scan(&message)
	}
	return message, nil
}

func (connection *MySqlStorageConnection) GetNextReceviedMessageToBeEnqueued() (*cap.CapReceivedMessage,error){
	statement := "SELECT * FROM `cap.received` WHERE `StatusName` = 'Scheduled' LIMIT 1;"
	conn, err:=connection.OpenDbConnection()
	if err != nil{
		return nil, err
	}
	rows,err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	message := &cap.CapReceivedMessage{}

	if rows.Next(){
		rows.Scan(&message)
	}

	return message,nil
}

func (connection *MySqlStorageConnection) GetPublishedMessage(id int) (*cap.CapPublishedMessage, error){
	statement := "SELECT * FROM `cap.published` WHERE `Id`=?;"
	conn, err:=connection.OpenDbConnection()
	if err != nil{
		return nil, err
	}
	rows,err := conn.Query(statement, id)
	if err != nil {
		return nil, err
	}
	message := &cap.CapPublishedMessage{}

	if rows.Next(){
		rows.Scan(&message)
	}

	return message,nil
}

func (connection *MySqlStorageConnection) GetReceivedMessage(id int) (*cap.CapReceivedMessage, error){
	statement := "SELECT * FROM `cap.received` WHERE `Id`=?;"
	conn, err:=connection.OpenDbConnection()
	if err != nil{
		return nil, err
	}
	rows,err := conn.Query(statement, id)
	if err != nil {
		return nil, err
	}
	message := &cap.CapReceivedMessage{}

	if rows.Next(){
		rows.Scan(&message)
	}

	return message,nil
}

func (connection *MySqlStorageConnection) StoreReceivedMessage(message *cap.CapReceivedMessage) error{
	statement := "INSERT INTO `{_prefix}.received`(`Name`,`Group`,`Content`,`Retries`,`Added`,`ExpiresAt`,`StatusName`)"
	statement += " VALUES(?,?,?,?,?,?,?);"
	conn, err:=connection.OpenDbConnection()
	if err != nil{
		return err
	}
	result ,err := conn.Exec(statement,message.Name,message.Group, message.Content, message.Retries, message.Added, message.ExpiresAt, message.StatusName)
	if err != nil{
		return err
	}
	rowAffected,err := result.RowsAffected()
	if err != nil{
		return err
	}
	if  rowAffected == int64(0) {
		return cap.NewCapError("Database execution should affect 1 row but affected 0 row actually.")
	}
	return nil
}