package cap

type InfiniteRetryProcessor struct{
	InnerProcessor IProcessor
}

func NewInfiniteRetryProcessor(innerProcessor IProcessor) *InfiniteRetryProcessor{
	processor := &InfiniteRetryProcessor{InnerProcessor:innerProcessor}
	return processor
}

func (this InfiniteRetryProcessor)Process(context *ProcessingContext){
	for{
		if context.IsStopping == false {
			err := this.InnerProcessor.Process(context)
			if err.Error() == "OperationCanceled" {
				return
			}else{
				//retry
			}
		}else{
			break
		}
	}
}