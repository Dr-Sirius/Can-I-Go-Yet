package src

import (
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/schedules"
	"can-i-go-yet/src/settings"
	"can-i-go-yet/src/templates"
)

var Schedules []schedules.Schedule

func Start() {
	
	settings.CreateSettings()
	handler.SetTime()
}
