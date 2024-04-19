package src

import (
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/handler"
)

var Schedules []scheduler.Schedule

func Start() {
	Schedules = scheduler.LoadSchedules()
	handler.SetTime()
}
