package ui

import (
	"can-i-go-yet/src/checker"
	"image/color"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())
	myWindow := app.NewWindow("Can I Go Yet?")
	todayLBL := canvas.NewText("Today's Schedule", color.White)
	todayLBL.TextSize = 35

	customerBTN := widget.NewButton("Customer View", func() {
		CustomerView()
	})

	content := container.New(
		layout.NewCenterLayout(),

		container.NewGridWithRows(
			2,
			todayLBL,
			DailyList(),
		),
		customerBTN,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(200, 200))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()

}

func DailyList() *widget.List {

	data := checker.GetSchedules()

	return widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {

			lbl := canvas.NewText("template", color.White)
			lbl.TextSize = 15
			return lbl
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = data[i].String()
		},
	)
}
