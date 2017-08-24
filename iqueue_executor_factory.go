package cap

const (
	PUBLISH = "Publish"
	SUBSCRIBE = "Subscribe"
)

type PublishQueueExecutorCreateDelegate func () IQueueExecutor

type IQueueExecutorFactory interface {
	SetPublishQueueExecutorCreateDelegate(delegate PublishQueueExecutorCreateDelegate)
	GetPublishQueueExecutorCreateDelegate() PublishQueueExecutorCreateDelegate
	GetInstance(messageType string) IQueueExecutor
}