package handler

import (
	"can-i-go-yet/src/converter"
	"can-i-go-yet/src/scheduler"
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2/container"
)

var sch []scheduler.Schedule
var stayOpen bool = false
var Announcments string = ""
var Status string = ""
var DefaultTemplate = ""

func SetTime() {
	s := scheduler.LoadSchedules()
	sch = scheduler.NewDayFromTime(time.Now(), s...).Schedules
	for _, x := range sch {
		if time.Now().Equal(x.StartTime) || (time.Now().After(x.StartTime) && time.Now().Before(x.EndTime)) {
			return
		}
	}
	if len(sch) == 0 && DefaultTemplate != "" {
		sch = converter.TemplateToSchedule(DefaultTemplate, time.Now().Format("2006-01-02"))
	}

}

func CheckTime() (string, color.Color) {
	if len(sch) == 0 {
		Status = "Open"
		return setOpen()

	} else if !GetCurrentSchedule().StartTime.Equal(time.Time{}) {
		return CheckFlags()
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
	flags := GetCurrentSchedule().Flags
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
	if checkSchedule() == -1 {
		return scheduler.Schedule{}
	}
	return sch[checkSchedule()]
}

func GetNextSchedule() scheduler.Schedule {
	if checkNextSchedule() == -1 {
		return scheduler.Schedule{}
	}
	return sch[checkNextSchedule()]
}

func checkNextSchedule() int {
	for i, x := range sch {
		if time.Now().Before(x.StartTime) {
			return i
		}
	}
	return -1
}

func checkSchedule() int {
	for i, x := range sch {
		if time.Now().Equal(x.StartTime) || (time.Now().After(x.StartTime) && time.Now().Before(x.EndTime)) {

			return i
		}
	}
	return -1
}

func GetStringSchedules() []string {
	s := []string{}
	for _, x := range sch {
		s = append(s, x.PrettyString())
	}
	return s
}

func GetReturnTime() string {
	return GetNextSchedule().StringStartTime()
}

func Remove(index int) {

	if (len(sch) == 1 || len(sch) == 0) && sch[0].Equal(scheduler.Schedule{}) {
		return
	}

	scheduler.RemoveSchedule(sch[index])
	s := sch[0:index]
	if len(sch)-1 > index {
		s = append(s, sch[index+1:]...)
	}

	if len(s) == 0 {
		s = append(s, scheduler.Schedule{})
	}
	sch = s

}

func RemoveAll() {
	os.WriteFile("Schedules/Schedules.csv", []byte("Date, Start_Time, End_Time, Flags"), os.ModePerm)
}

// func remove(s []scheduler.Schedule, index int) []scheduler.Schedule {
// 	sd := s[0:index]
// 	sd = append(sd, s[index+1:]...)
// 	return sd
// }

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

func SortTabs(tabs []*container.TabItem) []*container.TabItem {
	for i := len(tabs) - 1; i >= 0; i -= 1 {

		for x := range i {

			if tabs[x].Text > tabs[x+1].Text {
				temp := tabs[x]
				tabs[x] = tabs[x+1]
				tabs[x+1] = temp
			}

		}
	}
	return tabs
}
