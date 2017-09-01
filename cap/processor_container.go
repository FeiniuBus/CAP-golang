package cap

// ProcessorContainer bla.
type ProcessorContainer struct {
	Processors []IProcessor
}

// NewProcessorContainer bla.
func NewProcessorContainer() *ProcessorContainer {
	container := &ProcessorContainer{}
	container.Processors = make([]IProcessor, 0)
	return container
}

// Register bla.
func (container *ProcessorContainer) Register(processor IProcessor) *ProcessorContainer {
	container.Processors = append(container.Processors, processor)
	return container
}

// GetProcessors bla.
func (container *ProcessorContainer) GetProcessors() []*InfiniteRetryProcessor {
	processors := make([]*InfiniteRetryProcessor, 0)
	for _, val := range container.Processors {
		processors = append(processors, NewInfiniteRetryProcessor(val))
	}
	return processors
}
