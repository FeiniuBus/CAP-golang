package cap

type EnqueuedState struct {
	IState
}

func NewEnqueuedState() *EnqueuedState {
	return &EnqueuedState{
		
	}
}

func (this *EnqueuedState)GetExpiresAfter() int32{
	return 0
}

func (this *EnqueuedState)GetName() string{
	return "Enqueued"
}

func (this *EnqueuedState)ApplyReceivedMessage(message *CapReceivedMessage, transaction IStorageTransaction) error {
	return transaction.EnqueueReceivedMessage(message)
}

func (this *EnqueuedState)ApplyPublishedMessage(message *CapPublishedMessage, transaction IStorageTransaction) error {
	return transaction.EnqueuePublishedMessage(message)
}
