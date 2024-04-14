package ui

import (
	"can-i-go-yet/src/checker"
	"can-i-go-yet/src/scheduler"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())
	myWindow := app.NewWindow("Can I Go Yet?")

	content := container.NewAppTabs(
		container.NewTabItem("Today", DailyTab()),
		container.NewTabItem("Add Schedule", AddForm()),
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

			lbl := canvas.NewText("template", color.Black)
			lbl.TextSize = 15
			return lbl
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = data[i].PrettyString()
		},
	)
}

func DailyTab() *fyne.Container {
	todayLBL := canvas.NewText("Today's Schedule", color.Black)
	todayLBL.TextSize = 35

	currentLBL := canvas.NewText("Current Schedule: "+checker.GetCurrentSchedule().PrettyString(), color.Black)

	customerBTN := widget.NewButton("Customer View", func() {
		CustomerView()
	})

	go func() {
		for range time.Tick(time.Second) {
			currentLBL.Text = "Current Schedule: " + checker.GetCurrentSchedule().PrettyString()
		}
	}()
	return container.NewGridWithRows(
		4,
		todayLBL,
		DailyList(),
		currentLBL,
		customerBTN,
	)

}

func AddForm() *widget.Form {
	dtEntry := widget.NewEntry()
	dtEntry.SetPlaceHolder("2024-01-01")
	stEntry := widget.NewEntry()
	stEntry.SetPlaceHolder("12:00 am")
	etEntry := widget.NewEntry()
	etEntry.SetPlaceHolder("12:00 pm")
	flags := widget.CheckGroup{
		Horizontal: true,
		Required:   true,
		Options: []string{
			"Open",
			"Break",
			"Understaffed",
			"Holiday",
		},
	}
	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Date", Widget: dtEntry},
			{Text: "Start Time", Widget: stEntry},
			{Text: "End Time", Widget: etEntry},
			{Text: "Flags", Widget: &flags},
		},
		OnSubmit: func() {
			scheduler.AddSchedule(dtEntry.Text, stEntry.Text, etEntry.Text, checker.CreateFlags(flags.Selected)...)
		},
	}
}
