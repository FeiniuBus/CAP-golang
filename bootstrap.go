package cap

import (
	"github.com/FeiniuBus/capgo"
)

type ConsumeHandlerCreateDelegate func () IConsumerHandler

type Bootstrapper struct {
	Servers							[]cap.IProcessServer
	ConsumeHandlerDelegate			ConsumeHandlerCreateDelegate 
}

func NewBootstrapper() *Bootstrapper {
	rtv := &Bootstrapper{
		Servers: make([]cap.IProcessServer, 0) ,
	}

	return rtv
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

func initServers(bootstrapper *Bootstrapper) {

}