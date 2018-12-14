package plugin

import (
	"log"
	"runtime/debug"
)

// Recover func to get panic detail
func Recover(plgName string) {
	if v := recover(); v != nil {
		err, ok := v.(error)
		if !ok {
			log.Printf("plugin.%s panic: %s\n", plgName, debug.Stack())
			return
		}
		log.Printf("plugin.%s panic error: %v\n stack %s", plgName, err, debug.Stack())
	}
}
