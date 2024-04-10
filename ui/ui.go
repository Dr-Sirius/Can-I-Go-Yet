package ui

import (
	"can-i-go-yet/src"
	"can-i-go-yet/src/scheduler"
	"image/color"

	"time"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	app := app.New()
	//app.Settings().SetTheme(&defaultTheme{})
	myWindow := app.NewWindow("Can I Go Yet?")
	todayLBL := canvas.NewText("Today's Schedule", color.White)
	todayLBL.TextSize = 35

	content := container.New(
		layout.NewAdaptiveGridLayout(4),
		layout.NewSpacer(),
		layout.NewSpacer(),
		container.NewGridWithRows(
			5,
			layout.NewSpacer(),
			layout.NewSpacer(),
			todayLBL,
			DailyList(),
		),
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
			lbl := canvas.NewText("template",color.White)
			lbl.TextSize = 15
			return lbl
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = data.Schedules[i].String()
		},
	)
}
