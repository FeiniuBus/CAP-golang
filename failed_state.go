package cap

type FailedState struct {
	IState
}

func NewFailedState() *FailedState {
	return &FailedState{
		
	}
}

func (this *FailedState)GetExpiresAfter() int32{
	return 15 * 24 * 60 * 60
}

func (this *FailedState)GetName() string{
	return "Failed"
}

func (this *FailedState)ApplyReceivedMessage(message *CapReceivedMessage, transaction IStorageTransaction) error {
	return nil
}

func (this *FailedState)ApplyPublishedMessage(message *CapPublishedMessage, transaction IStorageTransaction) error {
	return nil
}
