package cap

type StorageConnectionFactory struct{
	CreateStorageConnection func(options CapOptions)(IStorageConnection, error)
}

func NewStorageConnectionFactory(delegate func(options CapOptions)(IStorageConnection,error)) *StorageConnectionFactory{
	factory := &StorageConnectionFactory{CreateStorageConnection : delegate}
	return factory
}