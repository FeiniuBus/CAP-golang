package cap

type IStateChanger interface {
	ChangeReceivedMessageState(message *CapReceivedMessage, state IState, transaction IStorageTransaction) error
	ChangePublishedMessage(message *CapPublishedMessage, state IState, transaction IStorageTransaction) error
}