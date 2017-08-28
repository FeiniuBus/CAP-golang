package cap

type StateChanger struct{
	
}

func NewStateChanger() IStateChanger{
	stateChanger := &StateChanger{}
	return stateChanger
}

func (this *StateChanger) ChangeReceivedMessageState(message *CapReceivedMessage, state IState, transaction IStorageTransaction) error{
	message.StatusName = state.GetName()
	err := state.ApplyReceivedMessage(message, transaction)
	if err != nil {
		return err
	}
	err = transaction.UpdateReceivedMessage(message)
	if err != nil {
		return err
	}
	return nil
}
func (this *StateChanger) ChangePublishedMessage(message *CapPublishedMessage, state IState, transaction IStorageTransaction) error{
	message.StatusName = state.GetName()
	err := state.ApplyPublishedMessage(message, transaction)
	if err != nil {
		return err
	}
	err = transaction.UpdatePublishedMessage(message)
	if err != nil {
		return err
	}
	return nil
}