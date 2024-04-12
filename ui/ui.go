package ui

import (
	"can-i-go-yet/src/checker"
	"image/color"

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

	content := container.NewDocTabs(
		container.NewTabItem("Today", DailyTab()),
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

func DailyTab() *container.Split {
	todayLBL := canvas.NewText("Today's Schedule", color.Black)
	todayLBL.TextSize = 35

	customerBTN := widget.NewButton("Customer View", func() {
		CustomerView()
	})
	return container.NewHSplit(
		container.NewGridWithRows(
			3,
			todayLBL,
			DailyList(),
			customerBTN,
		),
		AddForm(),
	)
}

func AddForm() *widget.Form {
	stEntry := widget.NewEntry()
	stEntry.SetPlaceHolder("12:00 am")
	etEntry := widget.NewEntry()
	etEntry.SetPlaceHolder("12:00 pm")
	flags := widget.CheckGroup{
		Horizontal: true,
		Options: []string{
			"Open",
			"Break",
			"Understaffed",
			"Holiday",
		},
	}
	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Start Time", Widget: stEntry},
			{Text: "End Time", Widget: etEntry},
			{Text: "Flags", Widget: &flags},
		},
	}
}
