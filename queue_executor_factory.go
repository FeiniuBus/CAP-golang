package cap

type QueueExecutorFactory struct {
	IQueueExecutorFactory
	PublishQueueExecutorCreateDelegate PublishQueueExecutorCreateDelegate
}

func (this *QueueExecutorFactory) SetPublishQueueExecutorCreateDelegate(delegate PublishQueueExecutorCreateDelegate){
	this.PublishQueueExecutorCreateDelegate = delegate
}

func (this *QueueExecutorFactory) GetPublishQueueExecutorCreateDelegate() PublishQueueExecutorCreateDelegate{
	return this.PublishQueueExecutorCreateDelegate
}

func (this *QueueExecutorFactory) GetInstance(messageType string) IQueueExecutor {
	if messageType == SUBSCRIBE {
		panic("todo impl")
	} else {
		return this.GetPublishQueueExecutorCreateDelegate()()
	}
}
