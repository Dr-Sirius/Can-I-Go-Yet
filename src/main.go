package src

import (
	"can-i-go-yet/src/scheduler"
	
)

var Schedules []scheduler.Schedule

func Start() {
	Schedules = scheduler.LoadSchedules()
}
