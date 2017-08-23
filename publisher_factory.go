package cap

type PublisherFactory struct{
	CreatePublisher func(options CapOptions)(IPublisher,error)
}

func NewPublisherFactory(createPublisher func(options CapOptions)(IPublisher,error)) *PublisherFactory{
	factory := &PublisherFactory{CreatePublisher:createPublisher}
	return factory
}