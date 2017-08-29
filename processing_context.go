package cap

import (
	"time"
)

// ProcessingContext bla.
type ProcessingContext struct {
	IsStopping bool
}

// NewProcessingContext bla.
func NewProcessingContext() *ProcessingContext {
	context := &ProcessingContext{IsStopping: false}
	return context
}

// ThrowIfStopping bla.
func (context *ProcessingContext) ThrowIfStopping() error {
	if context.IsStopping {
		return NewCapError("OperationCanceled")
	}
	return nil
}

// Stop bla.
func (context *ProcessingContext) Stop() {
	context.IsStopping = true
}

// Start bla.
func (context *ProcessingContext) Start() {
	context.IsStopping = false
}

// Wait bla.
func (context *ProcessingContext) Wait(d time.Duration) {
	time.Sleep(d)
}
