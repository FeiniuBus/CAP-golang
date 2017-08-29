package cap

import (
	"time"
)

// ProcessResult bla.
type ProcessResult struct {
	Status       string
	PollingDelay time.Duration
}

// ProcessSleeping bla.
func ProcessSleeping(pollingDelay time.Duration) *ProcessResult {
	return &ProcessResult{Status: "Sleeping", PollingDelay: pollingDelay}
}

// ProcessStoped bla.
func ProcessStoped() *ProcessResult {
	return &ProcessResult{Status: "Stop"}
}
