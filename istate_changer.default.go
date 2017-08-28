package cap

type StateChanger struct{
	
}

func NewStateChanger() IStateChanger{
	stateChanger := &StateChanger{}
	return stateChanger
}

func (this *StateChanger) ChangeReceivedMessageState(message *CapReceivedMessage, state IState, transaction IStorageTransaction) error{
	message.StatusName = state.GetName()
	err := state.ApplyReceivedMessage(message)
	if err != nil {
		return err
	}
	err := transaction.UpdateReceivedMessage(transaction, message)
	if err != nil {
		return err
	}
}
func (this *StateChanger) ChangePublishedMessage(message *CapPublishedMessage, state IState, transaction IStorageTransaction) error{
	message.StatusName = state.GetName()
	err := state.ApplyReceivedMessage(message)
	if err != nil {
		return err
	}
	err := transaction.UpdatePublishedMessage(transaction, message)
	if err != nil {
		return err
	}
}