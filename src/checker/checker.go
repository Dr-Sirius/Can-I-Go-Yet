package checker

import (
	"can-i-go-yet/src/scheduler"
	"image/color"
	"log"
	"time"
)

var sch []scheduler.Schedule
var currentSch int
var stayOpen bool = false

func SetTime() {
	s := scheduler.LoadSchedules()
	sch = scheduler.NewDayFromTime(time.Now(),s...).Schedules
	log.Println(sch)
	for i, x := range sch {
		if time.Now().Equal(x.StartTime) || (time.Now().After(x.StartTime) && time.Now().Before(x.EndTime)) {
			currentSch = i
			break
		}
	}
}

func CheckTime() (string, color.Color) {
	if len(sch) == 0 {
		log.Println("SCH = 0")
		return setOpen()

	} else if time.Now().Equal(sch[currentSch].StartTime) || time.Now().After(sch[currentSch].StartTime) {
		log.Println("Starttime")
		return CheckFlags()

	} else if time.Now().Equal(sch[currentSch].EndTime) || time.Now().After(sch[currentSch].EndTime) {
		log.Println("Endtime")
		if (currentSch + 1) != len(sch) {

			currentSch += 1
			CheckTime()
		} else {
			return setClosed()

		}

	} 
		
	if stayOpen {
		return setOpen()
	} else {
		return setClosed()
	}

}

func CheckFlags() (string, color.Color){
	flags := sch[currentSch].Flags
	if _, ok := flags[scheduler.BRKE]; ok {
		if _, ok := flags[scheduler.UNDS]; ok {
			return setClosed()
		} else {
			return setOnBreak()
		}

	} else if _, ok := flags[scheduler.OPEN]; ok {
		return setOpen()
	} else {
		return setClosed()
	}

}

func setOpen() (string, color.Color) {
	return "Open", color.RGBA{0, 255, 0, 255}
}

func setClosed() (string, color.Color) {
	return "Closed", color.RGBA{R: 255, G: 0, B: 0, A: 255}
}

func setOnBreak() (string, color.Color) {
	return "Open", color.RGBA{255, 255, 0, 255}
}


func GetSchedules() []scheduler.Schedule{
	return sch
}