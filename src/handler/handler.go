package handler

import (
	"can-i-go-yet/src/schedules"
	"can-i-go-yet/src/settings"
	"can-i-go-yet/src/templates"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2/container"
)

var sch []schedules.Schedule
var stayOpen bool = settings.LoadSettings().StayOpen
var Announcments string = ""
var Status string = ""
var defaultTemplate = settings.LoadSettings().DefaultTemplate

/*
Loads and Sets schedules for the current date
*/
func SetTime() {
	s, _ := templates.LoadTemplate(defaultTemplate)
	sch = s.Schedules
	for _, x := range sch {
		if time.Now().Equal(x.StartTime) || (time.Now().After(x.StartTime) && time.Now().Before(x.EndTime)) {
			return
		}
	}
	// if len(sch) == 0 && GetDefaultTemplate() != "" {
	// 	sch = converter.TemplateToSchedule(GetDefaultTemplate(), time.Now().Format("2006-01-02"))
	// }

}

/*
Checks and returns status of current schedule if there is one else returns status based on stayOpen var
*/
func CheckTime() (string, color.Color) {
	if len(sch) == 0 {
		goto stOpen

	} else if checkSchedule() != -1 {

		return CheckFlags()
	}
stOpen:
	if stayOpen {

		Status = "Open"
		return setOpen()
	} else {
		set := settings.LoadSettings()
		dHours := schedules.New(set.StandardHours[0], set.StandardHours[1], time.Now().Format("2006-01-02"), []int{0})
		if dHours.StartTime.Equal(time.Now()) || (dHours.StartTime.Before(time.Now()) && dHours.EndTime.After(time.Now())) {
			Status = "Open"
			return setOpen()
		}
		Status = "Closed"
		return setClosed()
	}

}

/*
Checks and returns status of current schedule based on its flags
*/
func CheckFlags() (string, color.Color) {
	flags := GetCurrentSchedule().Flags

	if schedules.HasFlag(flags, schedules.BRKE) {
		if schedules.HasFlag(flags, schedules.UNDS) {
			Status = "Closed"
			return setClosed()
		} else {
			Status = "Break"
			return setOnBreak()
		}

	} else if schedules.HasFlag(flags, schedules.UNDS) {
		return setUnderStaffed()
	} else if schedules.HasFlag(flags, schedules.OPEN) {
		if schedules.HasFlag(flags, -1) {
			Status = "Closed"
			return setClosed()
		}
		Status = "Open"
		return setOpen()
	} else {
		log.Println("closed")
		Status = "Closed"
		return setClosed()
	}

}

/*
Returns color.Color based on OpenColor settings in Settings.json
*/
func setOpen() (string, color.Color) {
	rgba := settings.LoadSettings().OpenColor
	return "Open", color.RGBA{uint8(rgba[0]), uint8(rgba[1]), uint8(rgba[2]), uint8(rgba[3])}
}

/*
Returns color.Color based on ClosedColor settings in Settings.json
*/
func setClosed() (string, color.Color) {
	rgba := settings.LoadSettings().ClosedColor
	return "Closed", color.RGBA{uint8(rgba[0]), uint8(rgba[1]), uint8(rgba[2]), uint8(rgba[3])}
}

/*
Returns color.Color based on BreakColor settings in Settings.json
*/
func setOnBreak() (string, color.Color) {
	rgba := settings.LoadSettings().BreakColor
	return "Open", color.RGBA{uint8(rgba[0]), uint8(rgba[1]), uint8(rgba[2]), uint8(rgba[3])}
}

func setUnderStaffed() (string, color.Color) {
	rgba := settings.LoadSettings().BreakColor
	return "Short Staffed", color.RGBA{uint8(rgba[0]), uint8(rgba[1]), uint8(rgba[2]), uint8(rgba[3])}
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
func GetSchedules() []schedules.Schedule {
	return sch
}

/*
Returns the current schedule from todays schedule based on current time
*/
func GetCurrentSchedule() schedules.Schedule {
	if checkSchedule() == -1 {
		return schedules.Schedule{}
	}
	return sch[checkSchedule()]
}

/*
Returns the next schedule from todays schedule based on current schedule and time
*/
func GetNextSchedule() schedules.Schedule {
	if checkNextSchedule() == -1 {
		return schedules.Schedule{StartTime: GetCurrentSchedule().EndTime}
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

	if (len(sch) == 1 || len(sch) == 0) && sch[0].Equal(schedules.Schedule{}) {
		return
	}

	//schedules.RemoveSchedule(sch[index])
	s := sch[0:index]
	if len(sch)-1 > index {
		s = append(s, sch[index+1:]...)
	}

	if len(s) == 0 {
		s = append(s, schedules.Schedule{})
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
Returns an array of ints containg schedules flag const values

Ex. []string{"Open","Understaffed","Break"} -> []int{0,2,1}
*/
func CreateFlags(flags []string) []int {
	f := []int{}
	if len(flags) == 0 {
		return []int{-1}
	}
	for _, x := range flags {
		if x == "Open" {
			f = append(f, schedules.OPEN)
		}
		if x == "Break" {
			f = append(f, schedules.BRKE)
		}
		if x == "Understaffed" {
			f = append(f, schedules.UNDS)
		}
		if x == "Holiday" {
			f = append(f, schedules.HDAY)
		}
	}
	log.Println(f)
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
