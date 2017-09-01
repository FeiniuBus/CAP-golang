package cap

type QueueExecutorFactory struct {
	IQueueExecutorFactory
	PublishQueueExecutorCreateDelegate PublishQueueExecutorCreateDelegate
	Register                           *CallbackRegister
}

func NewQueueExecutorFactory(register *CallbackRegister) *QueueExecutorFactory {
	return &QueueExecutorFactory{
		Register: register,
	}
}

func (this *QueueExecutorFactory) SetPublishQueueExecutorCreateDelegate(delegate PublishQueueExecutorCreateDelegate) {
	this.PublishQueueExecutorCreateDelegate = delegate
}

func (this *QueueExecutorFactory) GetPublishQueueExecutorCreateDelegate() PublishQueueExecutorCreateDelegate {
	return this.PublishQueueExecutorCreateDelegate
}

func (this *QueueExecutorFactory) GetInstance(messageType string) IQueueExecutor {
	if messageType == SUBSCRIBE {
		return NewQueueExecutorSubscribe(this.Register)
	} else {
		return this.PublishQueueExecutorCreateDelegate()
	}
}
