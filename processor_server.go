package cap

import (
	"sync"
	"time"
)

// ProcessorServer bla.
type ProcessorServer struct {
	Container  *ProcessorContainer
	Context    *ProcessingContext
	Processors []*InfiniteRetryProcessor
}

// NewProcessorServer bla.
func NewProcessorServer() *ProcessorServer {
	server := &ProcessorServer{Container: NewProcessorContainer(), Context: NewProcessingContext(), Processors: make([]*InfiniteRetryProcessor, 0)}
	return server
}

// Start bla.
func (server *ProcessorServer) Start() {
	server.Processors = server.Container.GetProcessors()
	for _, val := range server.Processors {
		go val.Process(server.Context)
	}
}

// StopTheWorld ...
func (server *ProcessorServer) StopTheWorld() chan bool {
	var result chan bool
	result <- true
	for _, val := range server.Processors {
		if val.Status == "Processing" {
			result <- false
		}
	}
	return result
}

// WaitForClose bla.
func (server *ProcessorServer) WaitForClose(wg *sync.WaitGroup) {
	server.Context.Stop()

	for {
		select {
		case <-server.StopTheWorld():
			wg.Done()
			//TODO : Kill process
			break

		default:
			time.Sleep(5 * time.Second)
		}
	}
}
