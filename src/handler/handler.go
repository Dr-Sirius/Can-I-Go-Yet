package handler

import (
	"can-i-go-yet/src/schedules"
	"can-i-go-yet/src/settings"
	"can-i-go-yet/src/templates"
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2/container"
)

var stayOpen bool = settings.LoadSettings().StayOpen
var Announcments string = ""
var Status string = ""
var currentSettings = settings.LoadSettings()
var defaultTemplate, defaultTemplateErr = templates.LoadTemplate(currentSettings.DefaultTemplate)
var currentSchedules []schedules.Schedule

func init() {
	if defaultTemplateErr != nil && currentSettings.DefaultTemplate != "" {
		log.Fatal(defaultTemplateErr)
	} else {
		defaultTemplate, _ = templates.LoadTemplate("TestTemplate")
	}
	currentSchedules = defaultTemplate.Schedules
}

/*
Returns the index of the schedule that StartTime is equal to time.Now() or the schedule that has time.Now() in between its StartTime and EndTime
*/
func getCurrentScheduleID() int {
	log.Println(currentSchedules)
	for i, x := range currentSchedules {
		log.Println(i, x)
		if x.StartTime.Equal(time.Now()) {
			return i
		} else if x.StartTime.Before(time.Now()) && x.EndTime.After(time.Now()) {
			return i
		}
	}
	return -1
}

func GetCurrentSchedule() schedules.Schedule {
	id := getCurrentScheduleID()
	if id == -1 {
		return schedules.Schedule{}
	}
	return currentSchedules[id]
}

func GetCurrentScheduleString() string {
	return GetCurrentSchedule().PrettyString()
}

func GetCurrentSchedules() []schedules.Schedule {
	return currentSchedules
}

func GetSchedule(ScheduleID int) schedules.Schedule {
	return currentSchedules[ScheduleID]
}

func GetReturnTime() string {
	return GetCurrentSchedule().StringStartTime()
}

func RemoveScheduleFromCurrent(schedule int) {
	currentSchedules = schedules.RemoveSchedule(schedule, currentSchedules)
}

func RemoveSchedulesFromCurrent() {
	currentSchedules = []schedules.Schedule{}
}

func CheckTime() (string, color.Color) {
	switch GetCurrentSchedule().Flags[0] {
	case schedules.OPEN:
		return setOpen()
	case schedules.BRKE:
		return setOnBreak()
	case schedules.UNDS:
		return setUnderStaffed()
	case schedules.HDAY:
		return setClosed()
	default:
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
	defaultTemplate, _ = templates.LoadTemplate(name)
}

func SetStayOpen(b bool) {
	stayOpen = b
}

func GetDefaultTemplate() string {
	return defaultTemplate.Name
}

func GetStayOpen() bool {
	return stayOpen
}

func Update() {
	defaultTemplate, _ = templates.LoadTemplate(settings.LoadSettings().DefaultTemplate)
	stayOpen = settings.LoadSettings().StayOpen
}
