package ui

import (
	"can-i-go-yet/src"
	"can-i-go-yet/src/scheduler"
	"log"
	"time"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)


func UserView(app fyne.App) {
	myWindow := app.NewWindow("Can I Go Yet?")
	myWindow.SetContent(DailyList())
	myWindow.ShowAndRun()

}

func DailyList() *widget.List {
	log.Println(src.Schedules)
	data := scheduler.NewDayFromTime(time.Now(), src.Schedules...)
	list := widget.NewList(
		func() int {
			return len(data.Schedules)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data.Schedules[i].String())
		},
	)
	return list
}
