package cap

import "fmt"

// UseConsoleLog ...
func UseConsoleLog(logger ILogger) {
	console := &LogDelegate{
		Log: func(level LogLevel, message string) {
			fmt.Println(string(level) + " : " + message)
		},
	}
}
