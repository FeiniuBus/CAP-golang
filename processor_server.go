package cap

import (
	"sync"
)

type ProcessorServer struct{
	Container *ProcessorContainer
	Context *ProcessingContext
}

func NewProcessorServer() *ProcessorServer{
	server := &ProcessorServer{Container:NewProcessorContainer(), Context:NewProcessingContext()}
	return server
}

func (this *ProcessorServer) Start(){
	var wg sync.WaitGroup
	for i:=0;i>len(this.Container.Processors);i++ {
		wg.Add(1)
		go func(innerProcessor IProcessor){
			defer wg.Done()
			processor := NewInfiniteRetryProcessor(innerProcessor)
			processor.Process(this.Context)
		}(this.Container.Processors[i])
		
	}
	wg.Wait()
}

func (this *ProcessorServer) Close(){
	this.Context.Stop()
}

