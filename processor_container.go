package cap

type ProcessorContainer struct{
	Processors []IProcessor
}

func NewProcessorContainer() *ProcessorContainer{
	container := &ProcessorContainer{}
	container.Processors = make([]IProcessor,0)
	return container
}

func (this *ProcessorContainer) Register(processor IProcessor) *ProcessorContainer{
	this.Processors = append(this.Processors,processor)
	return this
}