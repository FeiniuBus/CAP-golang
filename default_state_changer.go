package cap

import (
	"time"
)

type StateChanger struct {
	IStateChanger
}


func (this *StateChanger) ChangeReceivedMessageState(message *CapReceivedMessage, state IState, connection IStorageConnection) error {
	if state.GetExpiresAfter() != 0 {
		message.ExpiresAt = int(time.Now().Unix()) + int(state.GetExpiresAfter())
	} else {
		message.ExpiresAt = 0
	}

	transaction, err := connection.CreateTransaction()
	if err != nil {
		return err
	}

	message.StatusName = state.GetName()
	err = state.ApplyReceivedMessage(message, transaction)
	if err != nil {
		return err
	}

	err = transaction.UpdateReceivedMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (this *StateChanger) ChangePublishedMessage(message *CapPublishedMessage, state IState, connection IStorageConnection) error {
	if state.GetExpiresAfter() != 0 {
		message.ExpiresAt = int(time.Now().Unix()) + int(state.GetExpiresAfter())
	} else {
		message.ExpiresAt = 0
	}

	transaction, err := connection.CreateTransaction()
	if err != nil {
		return err
	}

	message.StatusName = state.GetName()
	err = state.ApplyPublishedMessage(message, transaction)
	if err != nil {
		return err
	}

	err = transaction.UpdatePublishedMessage(message)
	if err != nil {
		return err
	}
	
	return nil
}
