package cap

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// UseConsoleLog ...
func UseConsoleLog(logger *LoggerFactory) {
	console := &LogDelegate{
		Log: func(level LogLevel, message string, data interface{}) {
			text := level.GetName() + ":::" + message

			if data != nil {
				key := reflect.TypeOf(data).String()
				j, err := json.Marshal(data)
				if err == nil {
					text += ":::\r\n" + key + ":{" + string(j) + "}"
				}
			}

			fmt.Println(text)
		},
	}
	logger.Register(console)
}
