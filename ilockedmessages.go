package cap

// ILockedMessages
type ILockedMessages interface {
	Prepare(statement string) (stmt interface{}, err error)
	Commit() error
	Rollback() error
	Dispose()
	ChangeStates(state IState) error
	GetMessages() []ILockedMessage
	GetMessageType() int32
}
