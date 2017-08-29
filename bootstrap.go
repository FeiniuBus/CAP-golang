package cap

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Bootstrapper provide CAP booting functions.
type Bootstrapper struct {
	Servers              []IProcessServer
	CapOptions           *CapOptions
	Register             *CallbackRegister
	ConnectionFactory    *StorageConnectionFactory
	QueueExecutorFactory IQueueExecutorFactory
	WaitGroup            *sync.WaitGroup
}

// NewBootstrapper implements to instantiate an instance of Bootstrapper.
func NewBootstrapper(
	capOptions *CapOptions,
	connectionFactory *StorageConnectionFactory,
) *Bootstrapper {

	rtv := &Bootstrapper{
		Servers:           make([]IProcessServer, 0),
		CapOptions:        capOptions,
		Register:          NewCallbackRegister(),
		ConnectionFactory: connectionFactory,
	}

	initBootstrapper(rtv)

	return rtv
}

// initBootstrapper initlize Bootstrapper.
func initBootstrapper(bootstrapper *Bootstrapper) {
	bootstrapper.QueueExecutorFactory = &QueueExecutorFactory{}
}

// Bootstrap start CAP servers.
func (bootstrapper *Bootstrapper) Bootstrap() {
	for _, server := range bootstrapper.Servers {
		server.Start()
	}
}

// Close all started servers.
func (bootstrapper *Bootstrapper) Close() {
	for _, server := range bootstrapper.Servers {
		bootstrapper.WaitGroup.Add(1)
		go server.WaitForClose(bootstrapper.WaitGroup)
	}
	bootstrapper.WaitGroup.Wait()
}

// WaitForTerminalSignal ...
func (bootstrapper *Bootstrapper) WaitForTerminalSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	for {
		select {
		case s := <-c:
			fmt.Println("get signal:", s)
			bootstrapper.Close()
			break
		default:
			time.Sleep(10 * time.Second)
		}
	}
}

// Route listeners.
func (bootstrapper *Bootstrapper) Route(group, name string, cb CallbackInterface) {
	bootstrapper.Register.Add(group, name, cb)
}
