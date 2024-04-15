package checker

import (
	"can-i-go-yet/src/scheduler"
	"fmt"
	"image/color"
	"log"
	"time"
)

var sch []scheduler.Schedule
var currentSch int
var stayOpen bool = false
var Announcments string = ""
var Status string = "Break"

func SetTime() {
	s := scheduler.LoadSchedules()
	sch = scheduler.NewDayFromTime(time.Now(), s...).Schedules
	log.Println(sch)
	for i, x := range sch {
		if time.Now().Equal(x.StartTime) || (time.Now().After(x.StartTime) && time.Now().Before(x.EndTime)) {
			currentSch = i
			return
		}
	}
	currentSch = 0
	if len(sch) == 0 {
		sch = []scheduler.Schedule{
			{},
		}
	}

}

func CheckTime() (string, color.Color) {
	if len(sch) == 0 {
		log.Println("SCH = 0")
		return setOpen()

	} else if time.Now().Equal(sch[currentSch].StartTime) || (time.Now().After(sch[currentSch].StartTime) && time.Now().Before(sch[currentSch].EndTime)) {
		log.Println("Starttime")
		return CheckFlags()

	} else if time.Now().Equal(sch[currentSch].EndTime) || time.Now().After(sch[currentSch].EndTime) {
		log.Println("Endtime")
		if (currentSch + 1) != len(sch) {

			currentSch += 1
			CheckTime()
		} else {
			Status = "Closed"
			return setClosed()

		}

	}

	if stayOpen {
		Status = "Open"
		return setOpen()
	} else {
		Status = "Closed"
		return setClosed()
	}

}

func CheckFlags() (string, color.Color) {
	flags := sch[currentSch].Flags
	if _, ok := flags[scheduler.BRKE]; ok {
		if _, ok := flags[scheduler.UNDS]; ok {
			Status = "Closed"
			return setClosed()
		} else {
			Status = "Break"
			return setOnBreak()
		}

	} else if _, ok := flags[scheduler.OPEN]; ok {
		Status = "Open"
		return setOpen()
	} else {
		Status = "Closed"
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
	return "Open", color.RGBA{255, 218, 28, 255}
}

func GetDate() string {
	y, m, d := time.Now().Date()
	mt := fmt.Sprint(int(m))
	if int(m) < 10 {
		mt = "0" + mt
	}

	dt := fmt.Sprint(y) + "-" + mt + "-" + fmt.Sprint(d)
	return dt
}

func GetSchedules() []scheduler.Schedule {
	return sch
}

func GetCurrentSchedule() scheduler.Schedule {
	return sch[currentSch]
}

func GetStringSchedules() []string {
	s := []string{}
	for _, x := range sch {
		s = append(s, x.PrettyString())
	}
	return s
}

func GetReturnTime() string {
	
	if sch[currentSch].StartTime.After(time.Now()) {
		log.Println(true)
		return sch[currentSch].StringStartTime()
	}
	return sch[currentSch].StringEndTime()
}

func Remove(index int) {
	
	if (len(sch) == 1 || len(sch) == 0) && sch[0].Equal(scheduler.Schedule{}) {
		return
	}


	scheduler.RemoveSchedule(sch[index])
	s := sch[0:index]
	s = append(s, sch[index+1:]...)
	if len(s) == 0 {
		s = append(s, scheduler.Schedule{})
	}
	sch = s


}

func remove(s []scheduler.Schedule, index int) []scheduler.Schedule {
	sd := s[0:index]
	sd = append(sd, s[index+1:]...)
	return sd
}

func CreateFlags(flags []string) []int {
	f := []int{}
	for _, x := range flags {
		if x == "Open" {
			f = append(f, scheduler.OPEN)
		}
		if x == "Break" {
			f = append(f, scheduler.BRKE)
		}
		if x == "Understaffed" {
			f = append(f, scheduler.UNDS)
		}
		if x == "Holiday" {
			f = append(f, scheduler.HDAY)
		}
	}
	return f
}
