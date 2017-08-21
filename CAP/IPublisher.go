package cap


type IPublisher interface {  
	Publish(name string, content string, connection interface{}, transaction interface{}) error  
}
