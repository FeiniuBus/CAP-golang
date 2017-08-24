package cap

type Bootstrapper struct {
	Servers							[]IProcessServer
	CapOptions						*CapOptions
	Register						*CallbackRegister
	ConnectionFactory				*StorageConnectionFactory
	QueueExecutorFactory			IQueueExecutorFactory
}

func NewBootstrapper(
	capOptions *CapOptions, 
	register *CallbackRegister,
	connectionFactory *StorageConnectionFactory,
	) *Bootstrapper {

	rtv := &Bootstrapper{
		Servers: make([]IProcessServer, 0) ,
		CapOptions: capOptions,
		Register: register,
		ConnectionFactory: connectionFactory,
	}

	initBootstrapper(rtv)

	return rtv
}

func initBootstrapper(bootstrapper *Bootstrapper) {
	bootstrapper.QueueExecutorFactory = &QueueExecutorFactory {

	}
}

func (this *Bootstrapper) Bootstrap() {
	for _, server := range this.Servers {
		server.Start()
	}
}

func (this *Bootstrapper) Close() {
	for _, server := range this.Servers {
		server.Close()
	}
}

func (this *Bootstrapper) Route(group, name string, cb CallbackInterface) {
	this.Register.Add(group, name, cb)
}