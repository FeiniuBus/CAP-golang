package cap

// ILockedMessage ...
type ILockedMessage interface {
	Prepare(statement string) (stmt interface{}, err error)
	Commit() error
	Rollback() error
	Dispose()
	GetMessage() interface{}
	GetMessageType() int32
	Enqueue() (AffectedRows int64, err error)
	ChangeState(state IState) error
}
