package cap

import (
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

func (server *ProcessorServer) StopTheWorld() chan bool {
	var result chan bool
	result = make(chan bool, 1)
	for _, val := range server.Processors {
		if val.Status == "Processing" {
			result = make(chan bool, 0)
		}
	}
	return result
}

// Close bla.
func (server *ProcessorServer) Close() {
	server.Context.Stop()

	go func() {
		for {
			select {
			case <-server.StopTheWorld():
				//TODO : Kill process
				break

			default:
				time.Sleep(5 * time.Second)
			}
		}
	}()
}
