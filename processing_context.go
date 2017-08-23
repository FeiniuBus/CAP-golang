package cap

import (
	"time"
)

type ProcessingContext struct{
	IsStopping bool
}



func (this *ProcessingContext) WaitAsync(timeout time.Duration){
	time.Sleep(timeout)
}

func (this *ProcessingContext) ThrowIfStopping() error {
	if this.IsStopping {
		return NewCapError("OperationCanceled")
	}
	return nil
}

func (this *ProcessingContext) Stop(){
	this.IsStopping = true
}

func (this *ProcessingContext) Start(){
	this.IsStopping = false
}