package cap

type ScheduledState struct {
	IState
}

func NewScheduledState() *ScheduledState {
	return &ScheduledState{
		
	}
}

func (this *ScheduledState)GetExpiresAfter() int32{
	return 0
}

func (this *ScheduledState)GetName() string{
	return "Scheduled"
}

func (this *ScheduledState)ApplyReceivedMessage(message *CapReceivedMessage, transaction IStorageTransaction) error {
	return nil
}

func (this *ScheduledState)ApplyPublishedMessage(message *CapPublishedMessage, transaction IStorageTransaction) error {
	return nil
}
