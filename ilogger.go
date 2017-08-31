package cap

// ILogger ...
type ILogger interface {
	Log(level LogLevel, message string)
}
