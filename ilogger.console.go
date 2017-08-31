package cap

import "fmt"

// UseConsoleLog ...
func UseConsoleLog(logger *LoggerFactory) {
	console := &LogDelegate{
		Log: func(level LogLevel, message string) {
			fmt.Println(string(level) + ":::" + message)
		},
	}
	logger.Register(console)
}
