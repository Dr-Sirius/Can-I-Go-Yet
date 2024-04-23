package src

import (
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/settings"
	"can-i-go-yet/src/templater"
)

var Schedules []scheduler.Schedule

func Start() {
	Schedules = scheduler.LoadSchedules()
	templater.CreateTemplates()
	settings.CreateSettings()
	handler.SetTime()
}
