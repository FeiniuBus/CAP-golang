package cap

type StorageConnectionFactory struct{
	CreateStorageConnection func()(IStorageConnection, error)
}

func NewStorageConnectionFactory(delegate func()(IStorageConnection,error)) *StorageConnectionFactory{
	factory := &StorageConnectionFactory{CreateStorageConnection : delegate}
	return factory
}