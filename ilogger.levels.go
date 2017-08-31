package cap

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfomation
	LevelWarn
	LevelError
)