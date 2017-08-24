package cap

const (
	PUBLISH = "Publish"
	SUBSCRIBE = "Subscribe"
)

type IQueueExecutorFactory interface {
	GetInstance(messageType string) IQueueExecutor
}