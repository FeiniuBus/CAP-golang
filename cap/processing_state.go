package cap

type ProcessingState struct {
	IState
}

func NewProcessingState() *ProcessingState {
	return &ProcessingState{

	}
}

func (this *ProcessingState)GetExpiresAfter() int32{
	return 0
}

func (this *ProcessingState)GetName() string{
	return "Processing"
}

func (this *ProcessingState)ApplyReceivedMessage(message *CapReceivedMessage, transaction IStorageTransaction) error {
	return nil
}

func (this *ProcessingState)ApplyPublishedMessage(message *CapPublishedMessage, transaction IStorageTransaction) error {
	return nil
}
