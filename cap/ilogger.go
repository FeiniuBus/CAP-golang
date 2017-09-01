package cap

// ILogger ...
type ILogger interface {
	Log(level LogLevel, message string)
	LogData(level LogLevel, message string, data interface{})
}
