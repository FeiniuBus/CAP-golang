package cap

import (
	"sync"
)

// DefaultLogger ...
type DefaultLogger struct {
	delegates []*LogDelegate
}

var (
	singleton *DefaultLogger
	lock      sync.Mutex
)
var mutext *sync.Mutex

// NewDefaultLogger ...
func NewDefaultLogger() ILogger {
	if singleton == nil {
		lock.Lock()
		if singleton == nil {
			singleton = &DefaultLogger{
				delegates: make([]*LogDelegate, 0),
			}
		}
		lock.Unlock()
	}
	return singleton
}

// Log ...
func (logger *DefaultLogger) Log(level LogLevel, message string) {
	for _, val := range logger.delegates {
		val.Log(level, message)
	}
}

// Register ...
func (logger *DefaultLogger) Register(delegate *LogDelegate) {
	logger.delegates = append(logger.delegates, delegate)
}
