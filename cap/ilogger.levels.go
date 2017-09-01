package cap

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfomation
	LevelWarn
	LevelError
)

// GetName ...
func (level LogLevel) GetName() string {
	if level == LevelTrace {
		return "Trace"
	} else if level == LevelDebug {
		return "Debug"
	} else if level == LevelInfomation {
		return "Info"
	} else if level == LevelWarn {
		return "Warn"
	} else if level == LevelError {
		return "Error"
	}
	return "Unknown"
}
