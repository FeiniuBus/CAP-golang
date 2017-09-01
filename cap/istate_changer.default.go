package cap

// StateChanger ..
type StateChanger struct {
}

// NewStateChanger ..
func NewStateChanger() IStateChanger {
	stateChanger := &StateChanger{}
	return stateChanger
}

// ChangeReceivedMessageState ...
func (stateChanger *StateChanger) ChangeReceivedMessageState(message *CapReceivedMessage, state IState, transaction IStorageTransaction) error {
	message.StatusName = state.GetName()
	err := state.ApplyReceivedMessage(message, transaction)
	if err != nil {
		return err
	}
	message.StatusName = state.GetName()
	message.ExpiresAt = int(state.GetExpiresAfter())
	err = transaction.UpdateReceivedMessage(message)
	if err != nil {
		return err
	}
	return nil
}

// ChangePublishedMessage ...
func (stateChanger *StateChanger) ChangePublishedMessage(message *CapPublishedMessage, state IState, transaction IStorageTransaction) error {
	message.StatusName = state.GetName()
	err := state.ApplyPublishedMessage(message, transaction)
	if err != nil {
		return err
	}
	message.StatusName = state.GetName()
	message.ExpiresAt = int(state.GetExpiresAfter())
	err = transaction.UpdatePublishedMessage(message)
	if err != nil {
		return err
	}
	return nil
}
