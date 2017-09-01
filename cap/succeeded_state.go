package cap

type SucceededState struct {
	IState
}

func NewSucceededState() *SucceededState {
	return &SucceededState{}
}

func (this *SucceededState) GetExpiresAfter() int32 {
	return 1 * 60 * 60
}

func (this *SucceededState) GetName() string {
	return "Succeeded"
}

func (this *SucceededState) ApplyReceivedMessage(message *CapReceivedMessage, transaction IStorageTransaction) error {
	return nil
}

func (this *SucceededState) ApplyPublishedMessage(message *CapPublishedMessage, transaction IStorageTransaction) error {
	return nil
}
