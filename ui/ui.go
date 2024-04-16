package ui

import (
	"can-i-go-yet/src/checker"
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/templater"
	"log"

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
		container.NewTabItem("Today", TodayTab()),
		container.NewTabItem("Announcments", Announcments()),
		container.NewTabItem("Add Schedule", AddForm()),
		container.NewTabItem("Remove Schedule", Remove()),
		container.NewTabItem("Templates", TemplatTab()),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(200, 200))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()

}

func TodayList() *widget.List {

	return widget.NewList(
		func() int {
			data := checker.GetSchedules()
			return len(data)
		},
		func() fyne.CanvasObject {

			lbl := canvas.NewText("template", color.Black)
			lbl.TextSize = 15
			return lbl
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			data := checker.GetSchedules()
			o.(*canvas.Text).Text = data[i].PrettyString()
		},
	)
}

func TodayTab() *fyne.Container {
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
		TodayList(),
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
			if dtEntry.Text == checker.GetDate() {
				checker.SetTime()
			}
		},
	}
}

func Remove() *fyne.Container {
	selected := -1
	lbl := canvas.NewText("", color.Black)
	rl := TodayList()
	removeBTN := widget.NewButton("Remove", func() {
		if selected == -1 {
			lbl.Text = "You need to select a schedule first!"
			lbl.Refresh()
			return
		}

		checker.Remove(selected)
		selected = -1
		rl.UnselectAll()
		rl.Refresh()
	})

	rl.OnSelected = func(id widget.ListItemID) {
		lbl.Text = "Remove " + checker.GetSchedules()[id].PrettyString() + " ?"
		selected = id

		lbl.Refresh()
	}

	return container.NewGridWithRows(
		3,
		rl,
		lbl,
		removeBTN,
	)

}

func Announcments() *widget.Form {
	anc := widget.NewMultiLineEntry()

	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Anouncments", Widget: anc},
		},
		OnSubmit: func() {
			checker.Announcments = anc.Text
		},
	}
}

func TemplatTab() *fyne.Container {
	todayLBL := canvas.NewText("Templates", color.Black)
	todayLBL.TextSize = 35
	tabs := container.NewDocTabs(TemplateTabs()...)
	dateENT := widget.NewEntry()

	addBTN := widget.NewButton("Add Template", func() {
		name := ""
		tabs.OnSelected = func(ti *container.TabItem) {
			log.Println(ti.Text)
			name = ti.Text
		}
		
		for _, x := range templater.LoadTemplate(name) {

			scheduler.AddSchedule(dateENT.Text, x.Start_Time, x.End_Time, x.FlagsSlice()...)
		}
	})

	return container.NewGridWithRows(
		3,
		todayLBL,
		tabs,
		container.NewGridWithRows(
			2,
			dateENT,
			addBTN,
		),
	)

}

func TemplateTabs() []*container.TabItem {
	var tabs []*container.TabItem
	for i, x := range templater.GetAllTemplates() {
		c := container.NewTabItem(
			i,
			TemplateList(x),
		)
		tabs = append(tabs, c)
	}
	return tabs
}

func TemplateList(t []templater.Template) *widget.List {
	return widget.NewList(
		func() int {
			return len(t)
		},
		func() fyne.CanvasObject {

			lbl := canvas.NewText("template", color.Black)
			lbl.TextSize = 15
			return lbl
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = t[i].PrettyString()
		},
	)
}

func TemplateForm() *widget.Form {
	tName := widget.NewEntry()
	tName.SetPlaceHolder("Alpha")
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
			{Text: "Template Name", Widget: tName},
			{Text: "Start Time", Widget: stEntry},
			{Text: "End Time", Widget: etEntry},
			{Text: "Flags", Widget: &flags},
		},
		OnSubmit: func() {

		},
	}
}
