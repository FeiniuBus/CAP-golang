package cap

type IStateChanger interface {
	ChangeReceivedMessageState(message *CapReceivedMessage, state IState, connection IStorageConnection) error
	ChangePublishedMessage(message *CapPublishedMessage, state IState, connection IStorageConnection) error
}