package cap

type IPublishDelegate interface {
	Publish(keyName, content string) error
}
