package cap

// InfiniteRetryProcessor bla.
type InfiniteRetryProcessor struct {
	InnerProcessor IProcessor
	Status         string
}

// NewInfiniteRetryProcessor bla.
func NewInfiniteRetryProcessor(innerProcessor IProcessor) *InfiniteRetryProcessor {
	processor := &InfiniteRetryProcessor{InnerProcessor: innerProcessor, Status: "Stop"}
	return processor
}

// Process bla.
func (processor InfiniteRetryProcessor) Process(context *ProcessingContext) {
	for {
		if context.IsStopping == false {
			processor.Status = "Processing"
			result, err := processor.InnerProcessor.Process(context)
			if err != nil && err.Error() == "OperationCanceled" {
				return
			}
			if result != nil {
				processor.Status = result.Status
				if result.Status == "Sleeping" {
					context.Wait(result.PollingDelay)
				}
			}
		} else {
			processor.Status = "Stop"
			break
		}
	}
}
