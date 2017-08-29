package cap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

// Close bla.
func (server *ProcessorServer) WaitForClose() {
	server.Context.Stop()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-server.StopTheWorld():
				signal.Stop(c)
				//TODO : Kill process
				break

			default:
				time.Sleep(5 * time.Second)
			}
		}
	}()

	for {
		s := <-c
		fmt.Println("Got signal:", s)
	}
}
