package handler

import (
	"can-i-go-yet/src/converter"
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/settings"
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2/container"
)

var sch []scheduler.Schedule 
var stayOpen bool = settings.LoadSettings().StayOpen 
var Announcments string = ""
var Status string = ""
var defaultTemplate = settings.LoadSettings().DefaultTemplate

/*
Loads and Sets schedules for the current date
*/
func SetTime() {
	s := scheduler.LoadSchedules()
	sch = scheduler.NewDayFromTime(time.Now(), s...).Schedules
	for _, x := range sch {
		if time.Now().Equal(x.StartTime) || (time.Now().After(x.StartTime) && time.Now().Before(x.EndTime)) {
			return
		}
	}
	if len(sch) == 0 && GetDefaultTemplate() != "" {
		sch = converter.TemplateToSchedule(GetDefaultTemplate(), time.Now().Format("2006-01-02"))
	} 

}

/*
Checks and returns status of current schedule if there is one else returns status based on stayOpen var
*/
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

/*
Checks and returns status of current schedule based on its flags
*/
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

	} else if _, ok := flags[scheduler.UNDS]; ok{
		return setOnBreak()
	} else if _, ok := flags[scheduler.OPEN]; ok {
		Status = "Open"
		return setOpen()
	} else {
		Status = "Closed"
		return setClosed()
	}

}

/*
Returns color.Color based on OpenColor settings in Settings.json
*/
func setOpen() (string, color.Color) {
	rgba := settings.LoadSettings().OpenColor
	return "Open", color.RGBA{uint8(rgba[0]),uint8(rgba[1]),uint8(rgba[2]),uint8(rgba[3])}
}

/*
Returns color.Color based on ClosedColor settings in Settings.json
*/
func setClosed() (string, color.Color) {
	rgba := settings.LoadSettings().ClosedColor
	return "Closed", color.RGBA{uint8(rgba[0]),uint8(rgba[1]),uint8(rgba[2]),uint8(rgba[3])}
}

/*
Returns color.Color based on BreakColor settings in Settings.json
*/
func setOnBreak() (string, color.Color) {
	rgba := settings.LoadSettings().BreakColor
	return "Open", color.RGBA{uint8(rgba[0]),uint8(rgba[1]),uint8(rgba[2]),uint8(rgba[3])}
}

/*
Returns todays date as a string formmated as yyyy-mm-dd
*/
func GetDate() string {
	y, m, d := time.Now().Date()
	mt := fmt.Sprint(int(m))
	if int(m) < 10 {
		mt = "0" + mt
	}

	dt := fmt.Sprint(y) + "-" + mt + "-" + fmt.Sprint(d)
	return dt
}

/*
Returns all of todays schedules
*/
func GetSchedules() []scheduler.Schedule {
	return sch
}

/*
Returns the current schedule from todays schedule based on current time
*/
func GetCurrentSchedule() scheduler.Schedule {
	if checkSchedule() == -1 {
		return scheduler.Schedule{}
	}
	return sch[checkSchedule()]
}

/*
Returns the next schedule from todays schedule based on current schedule and time
*/
func GetNextSchedule() scheduler.Schedule {
	if checkNextSchedule() == -1 {
		return scheduler.Schedule{}
	}
	return sch[checkNextSchedule()]
}

/*
Returns index of next schedule based on current time
*/
func checkNextSchedule() int {
	for i, x := range sch {
		if time.Now().Before(x.StartTime) {
			return i
		}
	}
	return -1
}

/*
Returns index of current schedule based on current time
*/
func checkSchedule() int {
	for i, x := range sch {
		if time.Now().Equal(x.StartTime) || (time.Now().After(x.StartTime) && time.Now().Before(x.EndTime)) {

			return i
		}
	}
	return -1
}

/*
Returns a string array containing string versions of todays schedules
*/
func GetStringSchedules() []string {
	s := []string{}
	for _, x := range sch {
		s = append(s, x.PrettyString())
	}
	return s
}

/*
Returns string version of GetNextSchedule StartTime
*/
func GetReturnTime() string {
	return GetNextSchedule().StringStartTime()
}

/*
Removes given schedule at index from sch and Schedules.csv
*/
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

/*
Removes all schedules from sch and Schedules.csv
*/
func RemoveAll() {
	os.WriteFile("Schedules/Schedules.csv", []byte("Date, Start_Time, End_Time, Flags"), os.ModePerm)
}

/*
Returns an array of ints containg Scheduler flag const values

Ex. []string{"Open","Understaffed","Break"} -> []int{0,2,1}
*/
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

/*
Returns sorted []*container.TabItem using bubble sort
*/
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


func SetDefaultTemplate(name string) {
	defaultTemplate = name
}

func SetStayOpen(b bool) {
	stayOpen = b
}

func GetDefaultTemplate() string {
	return defaultTemplate
}

func GetStayOpen() bool {
	return stayOpen
}

func Update() {
	defaultTemplate = settings.LoadSettings().DefaultTemplate
	stayOpen = settings.LoadSettings().StayOpen
}