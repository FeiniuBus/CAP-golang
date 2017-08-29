package cap

// Bootstrapper provide CAP booting functions.
type Bootstrapper struct {
	Servers              []IProcessServer
	CapOptions           *CapOptions
	Register             *CallbackRegister
	ConnectionFactory    *StorageConnectionFactory
	QueueExecutorFactory IQueueExecutorFactory
}

// NewBootstrapper implements to instantiate an instance of Bootstrapper.
func NewBootstrapper(
	capOptions *CapOptions,
	register *CallbackRegister,
	connectionFactory *StorageConnectionFactory,
) *Bootstrapper {

	rtv := &Bootstrapper{
		Servers:           make([]IProcessServer, 0),
		CapOptions:        capOptions,
		Register:          register,
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
		server.Close()
	}
}

// Route listeners.
func (bootstrapper *Bootstrapper) Route(group, name string, cb CallbackInterface) {
	bootstrapper.Register.Add(group, name, cb)
}
