package cap

import (
	"time"
)

type CapOptions struct{
	ConnectionString string
	PoolingDelay time.Duration
	LogFunction func(message string,level ILogLevel)
}

func NewCapOptions() *CapOptions{
	options := &CapOptions{}
	return options
}

func (capOptions *CapOptions) GetConnectionString() (string, error){
	return capOptions.ConnectionString,nil	
}