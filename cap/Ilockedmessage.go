package cap

// ILockedMessage ...
type ILockedMessage interface {
	Prepare(statement string) (stmt interface{}, err error)
	Commit() error
	Rollback() error
	Dispose()
	GetMessage() interface{}
	GetMessageType() int32
	ChangeState(state IState) error
}
