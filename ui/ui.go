package ui

import (
	"can-i-go-yet/src/converter"
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/settings"
	"can-i-go-yet/src/templater"

	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*
Creates and Runs a new fyne.App - handles all ui related events
*/
func Run() {
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())

	myWindow := app.NewWindow("Can I Go Yet?")

	b := binding.NewUntypedList()
	for _, x := range scheduler.LoadSchedules() {
		b.Append(x)
	}

	content := container.NewAppTabs(
		container.NewTabItem("Today", TodayTab(b)),
		container.NewTabItem("Announcments", Announcments()),
		container.NewTabItem("Add Schedule", AddForm(b)),
		container.NewTabItem("Remove Schedule", Remove(b)),
		container.NewTabItem("Templates", TemplateTab(b)),
		container.NewTabItem("Build Template", BuildTemplatTab()),
		container.NewTabItem("Settings", SettingsTab(myWindow)),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(200, 200))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()

}

//region Schedules UI

/* 
Creates a *widget.List that contains todays schedules
*/
func TodayList(data binding.UntypedList) *widget.List {

	return widget.NewListWithData(
		data,
		func() fyne.CanvasObject {

			lbl := canvas.NewText("template", color.Black)
			lbl.TextSize = 15
			return lbl
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = converter.DataItemToSchedule(di).PrettyString()
		},
	)
}

/*
Creates a *fyne.Container tab that contains information about the todays date including customer view
*/
func TodayTab(data binding.UntypedList) *fyne.Container {
	todayLBL := canvas.NewText("Today's Schedule", color.Black)
	todayLBL.TextSize = 35

	currentLBL := canvas.NewText("Current Schedule: "+handler.GetCurrentSchedule().PrettyString(), color.Black)

	customerBTN := widget.NewButton("Customer View", func() {
		CustomerView()
	})

	go func() {
		for range time.Tick(time.Second) {
			currentLBL.Text = "Current Schedule: " + handler.GetCurrentSchedule().PrettyString()
		}
	}()
	return container.NewGridWithRows(
		4,
		todayLBL,
		TodayList(data),
		currentLBL,
		customerBTN,
	)

}

/*
Creates a *widget.Form for creating new schedule
*/
func AddForm(data binding.UntypedList) *widget.Form {
	dtEntry := widget.NewEntry()
	dtEntry.SetText(time.Now().Format("2006-01-02"))
	stEntry := widget.NewEntry()
	stEntry.SetText("12:00 am")
	etEntry := widget.NewEntry()
	etEntry.SetText("12:00 pm")
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
			s := scheduler.NewSchedule(stEntry.Text, etEntry.Text, dtEntry.Text, handler.CreateFlags(flags.Selected)...)
			scheduler.AddSchedule(dtEntry.Text, stEntry.Text, etEntry.Text, handler.CreateFlags(flags.Selected)...)
			data.Append(s)
			if dtEntry.Text == handler.GetDate() {
				handler.SetTime()
			}
		},
	}
}

/*
Creates *fyne.Container for removing schedules on todays date
*/
func Remove(data binding.UntypedList) *fyne.Container {
	selected := -1
	lbl := canvas.NewText("", color.Black)
	rl := TodayList(data)
	removeBTN := widget.NewButton("Remove", func() {
		if selected == -1 {
			lbl.Text = "You need to select a schedule first!"
			lbl.Refresh()
			return
		}

		handler.Remove(selected)
		s, _ := data.Get()
		temp := s[selected+1:]
		s = append(s[:selected], temp...)
		data.Set(s)
		selected = -1
		rl.UnselectAll()
		rl.Refresh()
	})

	removeAllBTN := widget.NewButton("Remove All", func() {
		handler.RemoveAll()
		s := make([]interface{}, 0)
		data.Set(s)
		rl.UnselectAll()
		rl.Refresh()
	})

	rl.OnSelected = func(id widget.ListItemID) {
		lbl.Text = "Remove " + handler.GetSchedules()[id].PrettyString() + " ?"
		selected = id

		lbl.Refresh()
	}

	return container.NewGridWithRows(
		3,
		rl,
		lbl,
		container.NewHBox(
			removeBTN,
			removeAllBTN,
		),
	)

}

/*
Creates *widget.Form for adding anouncment to customer view
*/
func Announcments() *widget.Form {
	anc := widget.NewMultiLineEntry()

	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Anouncments", Widget: anc},
		},
		OnSubmit: func() {
			handler.Announcments = anc.Text
		},
	}
}

//region Templates UI

/*
Creates *fyne.Container for viewing templates and adding them to todays schedule
*/
func TemplateTab(data binding.UntypedList) *fyne.Container {
	todayLBL := canvas.NewText("Templates", color.Black)
	todayLBL.TextSize = 35
	tabs := container.NewDocTabs(TemplateTabs()...)
	dateENT := widget.NewEntry()
	dateENT.SetText(time.Now().Format("2006-01-02"))
	name := ""
	if len(tabs.Items) != 0 {
		name = tabs.Items[0].Text
	}

	tabs.OnSelected = func(ti *container.TabItem) {
		name = ti.Text
	}

	tabs.OnClosed = func(ti *container.TabItem) {
		templater.RemoveTemplate(ti.Text)

	}

	addBTN := widget.NewButton("Add Template", func() {

		for _, x := range templater.LoadTemplate(name) {
			
			scheduler.AddSchedule(dateENT.Text, x.Start_Time, x.End_Time, x.FlagsSlice()...)
			data.Append(scheduler.NewSchedule(x.Start_Time,x.End_Time,dateENT.Text, x.FlagsSlice()...))
			
		}
		
		handler.SetTime()

	})

	go func() {
		for range time.Tick(time.Second) {
			tabs.SetItems(TemplateTabs())
		}
	}()

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


/*
Creates []*container.TabItem that holds tabs for each saved template in Templates folder
*/
func TemplateTabs() []*container.TabItem {
	var tabs []*container.TabItem
	t := templater.GetAllTemplates()
	if len(t) != 0 {
		for i, x := range  t{
			c := container.NewTabItem(
				i,
				TemplateList(x),
			)
	
			tabs = append(tabs, c)
		}
	}
	

	return handler.SortTabs(tabs)
}


/*
Creates *widget.List that displays template information for the passed []templater.Template
*/
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

//region Build Templates UI

/*
Creates *widget.Form for creating new template
*/
func TemplateForm(list *widget.List, b *binding.UntypedList) *widget.Form {
	tName := widget.NewEntry()
	
	stEntry := widget.NewEntry()
	stEntry.SetPlaceHolder("12:00 am")
	etEntry := widget.NewEntry()
	etEntry.SetPlaceHolder("12:00 pm")
	saveBTN := widget.NewButton("Save", func() {
		t := []templater.Template{}
		for i := range (*b).Length() {
			item, _ := (*b).GetItem(i)
			t = append(t, converter.DataItemToTemplate(item))
		}
		templater.AddTemplate(t)
		tName.Text = ""
	})
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
			{Widget: saveBTN},
		},
		OnSubmit: func() {
			(*b).Append(templater.NewTemplate(tName.Text, stEntry.Text, etEntry.Text, handler.CreateFlags(flags.Selected)...))

			stEntry.Text = ""
			etEntry.Text = ""
			stEntry.Refresh()
			etEntry.Refresh()
			list.Refresh()
		},
	}
}


/*
Creates *widget.List for displaying information about template being made in BuildTemplateTab()
*/
func BuildTemplateList(data binding.UntypedList) *widget.List {
	return widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			lbl := canvas.NewText("template", color.Black)
			lbl.TextSize = 15
			return lbl
		},
		func(di binding.DataItem, o fyne.CanvasObject) {

			o.(*canvas.Text).Text = converter.DataItemToTemplate(di).PrettyString()
		},
	)
}

/*
Creates *container.Split for building new templates
*/
func BuildTemplatTab() *container.Split {

	b := binding.NewUntypedList()

	ls := BuildTemplateList(b)

	content := container.NewVSplit(
		ls,
		TemplateForm(ls, &b),
	)

	return content
}

//region Settings UI


/*
Creates *widget.Form for displaying and changing settings - reads and writes to Settings/Settings.json
*/
func SettingsTab(w fyne.Window) *widget.Form {
	tName := widget.NewEntry()
	tName.Text = handler.GetDefaultTemplate()
	stoCheck := widget.NewCheck("", func(b bool) { handler.SetStayOpen(b) })

	var oColor color.Color = converter.IntToColor(settings.LoadSettings().OpenColor)
	opRect := canvas.NewRectangle(oColor)
	var bColor color.Color = converter.IntToColor(settings.LoadSettings().BreakColor)
	bkRect := canvas.NewRectangle(bColor)
	var cColor color.Color = converter.IntToColor(settings.LoadSettings().ClosedColor)
	clRect := canvas.NewRectangle(cColor)

	openColorDialog := dialog.NewColorPicker("Open Color", "", func(c color.Color) { oColor = c; SetColor(c, opRect) }, w)
	openColorDialog.Advanced = true
	closedColorDialog := dialog.NewColorPicker("Closed Color", "", func(c color.Color) { cColor = c; SetColor(c, clRect) }, w)
	closedColorDialog.Advanced = true
	breakColorDialog := dialog.NewColorPicker("Break Color", "", func(c color.Color) { bColor = c; SetColor(c, bkRect) }, w)
	breakColorDialog.Advanced = true

	opBTN := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), openColorDialog.Show)
	clBTN := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), closedColorDialog.Show)
	bkBTN := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), breakColorDialog.Show)

	oContent := container.NewGridWithColumns(2, opRect, opBTN)
	cContent := container.NewGridWithColumns(2, clRect, clBTN)
	bContent := container.NewGridWithColumns(2, bkRect, bkBTN)

	stEntry := widget.NewEntry()
	stEntry.SetPlaceHolder("7:30 am")

	etEntry := widget.NewEntry()
	etEntry.SetPlaceHolder("2:30 pm")

	dhContent := container.NewVBox(stEntry,etEntry)

	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Default Template", Widget: tName},
			{Text: "Stay Open", Widget: stoCheck},
			{Text: "Open Label Color", Widget: oContent},
			{Text: "Closed Label Color", Widget: cContent},
			{Text: "Break Label Color", Widget: bContent},
			{Text: "Daily Hours",Widget: dhContent},
		},
		SubmitText: "Save",
		OnSubmit: func() {

			s := settings.Settings{
				DefaultTemplate: tName.Text,
				StayOpen:        stoCheck.Checked,
				OpenColor:       converter.ColorToInt(oColor),
				ClosedColor:     converter.ColorToInt(cColor),
				BreakColor:      converter.ColorToInt(bColor),
				StandardHours:   [2]string{stEntry.Text, etEntry.Text},
			}
			s.SaveSettings()
			handler.Update()

		},
	}
}

/*
Sets passed *canvas.Rectange color to passed color.Color
*/
func SetColor(c color.Color, r *canvas.Rectangle) {
	r.FillColor = c
}
