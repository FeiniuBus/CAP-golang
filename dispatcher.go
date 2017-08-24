package cap

import (
	"time"
)

type Dispatcher struct{
	Options CapOptions
	StorageConnectionFactory *StorageConnectionFactory
	PublishQueueProccesor *PublishQueuer
}

func NewDispatcher(options CapOptions, storageConnectionFactory *StorageConnectionFactory) *Dispatcher{
	 dispatcher := &Dispatcher{StorageConnectionFactory:storageConnectionFactory}	
	 return dispatcher
}

func (this Dispatcher) ExecutePublishQueuer(){
	storageConnection, err := this.StorageConnectionFactory.CreateStorageConnection(this.Options)
	if err != nil{
		panic(err)
	}
	queueProcessor := NewPublishQueuer(storageConnection)
	err = queueProcessor.Execute()
	if err != nil{
		panic(err)
	}
}

func (this Dispatcher) Begin(){
	tick := time.Tick(60 * time.Second)
	for{
		select{
		case <- tick :
			go this.ExecutePublishQueuer()
		default:
			time.Sleep(60 * time.Second)
		}
	}

}