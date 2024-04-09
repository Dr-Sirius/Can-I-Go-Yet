package ui

import (
	"can-i-go-yet/src"
	"can-i-go-yet/src/scheduler"
	"time"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	app := app.New()
	//app.Settings().SetTheme(&defaultTheme{})
	myWindow := app.NewWindow("Can I Go Yet?")

	content := container.New(
		layout.NewAdaptiveGridLayout(3),
		widget.NewLabel("Today's Schedules"),
		DailyList(),
	)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func DailyList() *widget.List {
	data := scheduler.NewDayFromTime(time.Now(), src.Schedules...)
	return widget.NewList(
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
}
