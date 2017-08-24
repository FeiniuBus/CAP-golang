package cap

type IQueueExecutor interface {
	Execute(connection IStorageConnection, feched IFetchedMessage) error
}