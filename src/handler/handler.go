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
var defaultTemplate, defaultTemplateErr = templates.LoadTemplate(settings.LoadSettings().DefaultTemplate)

func init() {
	if defaultTemplateErr != nil {
		log.Fatal(defaultTemplateErr)
	}
}

var currentSchedules = defaultTemplate.Schedules

/*
Returns the index of the schedule that StartTime is equal to time.Now() or the schedule that has time.Now() in between its StartTime and EndTime
*/
func getCurrentScheduleID() int {
	for i, x := range currentSchedules {
		if x.StartTime.Equal(time.Now()) {
			return i
		} else if x.StartTime.Before(time.Now()) && x.EndTime.After(time.Now()) {
			return i
		}
	}
	return -1
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
