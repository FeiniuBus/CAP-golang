package cap

type ILockedMessage interface {
	Prepare(statement string) (stmt interface{}, err error)
	Commit() error
	Rollback() error
	Dispose()
	GetMessage() interface{}
	GetMessageType() int32
}
