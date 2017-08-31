package cap

// LogDelegate ...
type LogDelegate struct {
	// Log ...
	Log func(level LogLevel, message string, data interface{})
}
