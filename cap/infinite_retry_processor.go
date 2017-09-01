package cap

// InfiniteRetryProcessor bla.
type InfiniteRetryProcessor struct {
	InnerProcessor IProcessor
	Status         string
	logger         ILogger
}

// NewInfiniteRetryProcessor bla.
func NewInfiniteRetryProcessor(innerProcessor IProcessor) *InfiniteRetryProcessor {
	processor := &InfiniteRetryProcessor{InnerProcessor: innerProcessor, Status: "Stop"}
	processor.logger = GetLoggerFactory().CreateLogger(processor)
	return processor
}

// Process bla.
func (processor InfiniteRetryProcessor) Process(context *ProcessingContext) {
	for {
		if context.IsStopping == false {
			processor.Status = "Processing"
			result, err := processor.InnerProcessor.Process(context)
			if err != nil {
				processor.logger.Log(LevelError, "[Process]"+err.Error())
			}
			if err != nil && err.Error() != "OperationCanceled" {
				processor.logger.Log(LevelError, "[Process]"+err.Error())
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
