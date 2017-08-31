package cap

import (
	"reflect"
	"sync"
)

// LoggerFactory ...
type LoggerFactory struct {
	delegates []*LogDelegate
}

var (
	loggerFactoryInstance *LoggerFactory
	loggerFactoryLocker   sync.Mutex
)

// GetLoggerFactory ...
func GetLoggerFactory() *LoggerFactory {
	if loggerFactoryInstance == nil {
		loggerFactoryLocker.Lock()
		if loggerFactoryInstance == nil {
			loggerFactoryInstance = &LoggerFactory{
				delegates: make([]*LogDelegate, 0),
			}
		}
		loggerFactoryLocker.Unlock()
	}
	return loggerFactoryInstance
}

// Register ...
func (factory *LoggerFactory) Register(delegate *LogDelegate) {
	factory.delegates = append(factory.delegates, delegate)
}

// CreateLogger ...
func (factory *LoggerFactory) CreateLogger(i interface{}) ILogger {
	t := reflect.TypeOf(i)
	logger := &DefaultLogger{
		TypeName:  t.String(),
		delegates: factory.delegates,
	}
	return logger
}
