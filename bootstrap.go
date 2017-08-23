package cap

type Bootstrapper struct {
	Servers							[]IProcessServer
	CapOptions						CapOptions
	Register						*CallbackRegister
	ConnectionFactory				*StorageConnectionFactory
}

func NewBootstrapper(
	capOptions CapOptions, 
	register *CallbackRegister,
	connectionFactory *StorageConnectionFactory,
	) *Bootstrapper {
		
	rtv := &Bootstrapper{
		Servers: make([]IProcessServer, 0) ,
		CapOptions: capOptions,
		Register: register,
		ConnectionFactory: connectionFactory,
	}

	initBootstrpper(rtv)

	return rtv
}

func initBootstrpper(bootstrapper *Bootstrapper) {

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
