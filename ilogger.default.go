package cap

import (
	"sync"
)

// DefaultLogger ...
type DefaultLogger struct {
	TypeName  string
	delegates []*LogDelegate
}

var mutext *sync.Mutex

// Log ...
func (logger *DefaultLogger) Log(level LogLevel, message string) {
	for _, val := range logger.delegates {
		val.Log(level, logger.TypeName+"->"+message, nil)
	}
}

// LogData ...
func (logger *DefaultLogger) LogData(level LogLevel, message string, data interface{}) {
	for _, val := range logger.delegates {
		val.Log(level, logger.TypeName+"->"+message, data)
	}
}
