package cap

// IPublisher ...
type IPublisher interface {
	Publish(descriptors []*MessageDescriptor, connection interface{}, transaction interface{}) error

	PublishOne(name string, content interface{}, connection interface{}, transaction interface{}) error
}
