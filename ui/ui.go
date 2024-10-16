package ui

import (
	"can-i-go-yet/src/converter"
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/schedules"
	"can-i-go-yet/src/settings"
	"can-i-go-yet/src/templates"
	"errors"

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
	// for _, x := range schedules.LoadSchedules() {
	// 	b.Append(x)
	// }

	content := container.NewAppTabs(
		container.NewTabItem("Today", TodayTab(b)),
		container.NewTabItem("Announcments", Announcments()),
		container.NewTabItem("Add Schedule", AddForm(b, myWindow)),
		container.NewTabItem("Remove Schedule", Remove(b, myWindow)),
		container.NewTabItem("Templates", TemplateTab(b, myWindow)),
		container.NewTabItem("Build Template", BuildTemplatTab(myWindow)),
		container.NewTabItem("Settings", SettingsTab(myWindow)),
	)

	myWindow.SetContent(content)
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
func AddForm(data binding.UntypedList, win fyne.Window) *widget.Form {
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
			"Closed",
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
			if dtEntry.Text == "" || stEntry.Text == "" || etEntry.Text == "" {

				dialog.NewError(errors.New("the entries cannot be left blank"), win).Show()

				return

			}
			s := schedules.New(stEntry.Text, etEntry.Text, dtEntry.Text, handler.CreateFlags(flags.Selected))
			//schedules.AddSchedule(dtEntry.Text, stEntry.Text, etEntry.Text, handler.CreateFlags(flags.Selected)...)
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
func Remove(data binding.UntypedList, win fyne.Window) *fyne.Container {
	selected := -1
	lbl := canvas.NewText("", color.Black)
	rl := TodayList(data)
	removeBTN := widget.NewButton("Remove", func() {
		if selected == -1 {

			dialog.NewError(errors.New("there are no schedules selected"), win).Show()

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

	// removeAllBTN := widget.NewButton("Remove All", func() {
	// 	if len(schedules.LoadSchedules()) == 0 {

	// 		dialog.NewError(errors.New("there are no schedules\n\n you can create new schedules in the Add Schedule tab"),win).Show()

	// 		return

	// 	}
	// 	handler.RemoveAll()
	// 	s := make([]interface{}, 0)
	// 	data.Set(s)
	// 	rl.UnselectAll()
	// 	rl.Refresh()
	// })

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
			//removeAllBTN,
		),
	)

}

/*
Creates *widget.Form for adding anouncment to customer view
*/
func Announcments() *widget.Form {
	anc := widget.NewEntry()

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
func TemplateTab(data binding.UntypedList, win fyne.Window) *fyne.Container {
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
		s := settings.LoadSettings()

		if s.ShowDeleteConfirmation {
			w := widget.NewCheck("Do not show again", func(b bool) {
				if b {
					s.ShowDeleteConfirmation = false
					s.SaveSettings()

				}
			})

			dialog.NewCustomConfirm(
				"This action will delete the selected template\nDo you wish to continue with this action?",
				"Yes",
				"No",
				w, func(b bool) {
					if b {
						templates.RemoveTemplate(ti.Text)
					}
				},
				win,
			).Show()
			return

		}
		templates.RemoveTemplate(ti.Text)

	}

	addBTN := widget.NewButton("Add Template to Todays Schedule", func() {
		if len(templates.LoadAllTemplates()) == 0 {
			dialog.NewError(errors.New("there are no templates\n\n you can create a new template in the Build Template tab"), win).Show()

			return
		}
		if name == "" {
			dialog.NewError(errors.New("there are no templates currently selected"), win).Show()

			return
		}
		template, _ := templates.LoadTemplate(name)
		for _, x := range template.Schedules {

			data.Append(schedules.New(x.StringStartTime(), x.StringEndTime(), dateENT.Text, x.Flags))

		}

		handler.SetTime()

	})

	replaceBTN := widget.NewButton("Replace Todays Schedule with Template", func() {
		if len(templates.LoadAllTemplates()) == 0 {
			dialog.NewError(errors.New("there are no templates\n\n you can create a new template in the Build Template tab"), win).Show()

			return
		}
		if name == "" {
			dialog.NewError(errors.New("there are no templates currently selected"), win).Show()

			return
		}
		s := make([]interface{}, 0)
		data.Set(s)
		handler.RemoveAll()
		template, _ := templates.LoadTemplate(name)
		for _, x := range template.Schedules {

			data.Append(schedules.New(x.StringStartTime(), x.StringEndTime(), dateENT.Text, x.Flags))

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
			container.NewHBox(
				addBTN,
				replaceBTN,
			),
		),
	)

}

/*
Creates []*container.TabItem that holds tabs for each saved template in Templates folder
*/
func TemplateTabs() []*container.TabItem {
	var tabs []*container.TabItem
	template := templates.LoadAllTemplates()
	if len(template) != 0 {
		for _, x := range template {
			c := container.NewTabItem(
				x.Name,
				TemplateList(x.Schedules),
			)

			tabs = append(tabs, c)
		}
	}

	return handler.SortTabs(tabs)
}

/*
Creates *widget.List that displays template information for the passed []templates.Template
*/
func TemplateList(s []schedules.Schedule) *widget.List {
	return widget.NewList(
		func() int {
			return len(s)
		},
		func() fyne.CanvasObject {

			lbl := canvas.NewText("template", color.Black)
			lbl.TextSize = 15
			return lbl
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = s[i].PrettyString()
		},
	)
}

//region Build Templates UI

/*
Creates *widget.Form for creating new template
*/
func TemplateForm(list *widget.List, b *binding.UntypedList, win fyne.Window) *fyne.Container {
	tName := widget.NewEntry()

	stEntry := widget.NewEntry()
	stEntry.SetText("12:00 am")
	etEntry := widget.NewEntry()
	etEntry.SetText("12:00 pm")
	saveBTN := widget.NewButtonWithIcon("Save Template", theme.DocumentSaveIcon(), func() {
		schedule := []schedules.Schedule{}
		for i := range (*b).Length() {
			item, _ := (*b).GetItem(i)
			schedule = append(schedule, converter.DataItemToSchedule(item))
		}
		templates.CreateTemplateFile(templates.Template{Name: tName.Text, Schedules: schedule})
		tName.Enable()
		s := make([]interface{}, 0)
		(*b).Set(s)
		stEntry.Text = "12:00 am"
		etEntry.Text = "12:00 pm"
		tName.Text = ""
		tName.Refresh()
		stEntry.Refresh()
		etEntry.Refresh()
		list.Refresh()

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

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Template Name", Widget: tName},
			{Text: "Start Time", Widget: stEntry},
			{Text: "End Time", Widget: etEntry},
			{Text: "Flags", Widget: &flags},
		},
		SubmitText: "Add",
		OnSubmit: func() {
			if tName.Text == "" || stEntry.Text == "" || etEntry.Text == "" {

				dialog.NewError(errors.New("the entries cannot be left blank"), win).Show()

				return

			}
			tName.Disable()
			(*b).Append(schedules.New(tName.Text, stEntry.Text, etEntry.Text, handler.CreateFlags(flags.Selected)))

			stEntry.Text = ""
			etEntry.Text = ""
			stEntry.Refresh()
			etEntry.Refresh()
			list.Refresh()
		},
		OnCancel: func() {
			tName.Enable()
			s := make([]interface{}, 0)
			(*b).Set(s)
			stEntry.Text = "12:00 am"
			etEntry.Text = "12:00 pm"
			tName.Text = ""
			tName.Refresh()
			stEntry.Refresh()
			etEntry.Refresh()
			list.Refresh()
		},
	}
	content := container.NewVBox(
		form,
		saveBTN,
	)
	return content
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

			o.(*canvas.Text).Text = converter.DataItemToSchedule(di).PrettyString()
		},
	)
}

/*
Creates *container.Split for building new templates
*/
func BuildTemplatTab(win fyne.Window) *container.Split {

	b := binding.NewUntypedList()

	ls := BuildTemplateList(b)

	content := container.NewVSplit(
		ls,
		TemplateForm(ls, &b, win),
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
	stEntry.SetText(settings.LoadSettings().StandardHours[0])

	etEntry := widget.NewEntry()
	etEntry.SetText(settings.LoadSettings().StandardHours[1])

	dhContent := container.NewVBox(stEntry, etEntry)

	fsCheck := widget.NewCheck("", func(b bool) {})
	fsCheck.SetChecked(settings.LoadSettings().FullscreenCustomerView)

	dnsaCheck := widget.NewCheck("", func(b bool) {})
	dnsaCheck.SetChecked(settings.LoadSettings().ShowDeleteConfirmation)

	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Default Template", Widget: tName},
			{Text: "Stay Open", Widget: stoCheck},
			{Text: "Open Label Color", Widget: oContent},
			{Text: "Closed Label Color", Widget: cContent},
			{Text: "Break Label Color", Widget: bContent},
			{Text: "Daily Hours", Widget: dhContent},
			{Text: "Fullscreen Customer View", Widget: fsCheck},
			{Text: "Show delete confirmation", Widget: dnsaCheck},
		},
		SubmitText: "Save",
		OnSubmit: func() {

			s := settings.Settings{
				DefaultTemplate:        tName.Text,
				StayOpen:               stoCheck.Checked,
				OpenColor:              converter.ColorToInt(oColor),
				ClosedColor:            converter.ColorToInt(cColor),
				BreakColor:             converter.ColorToInt(bColor),
				StandardHours:          [2]string{stEntry.Text, etEntry.Text},
				FullscreenCustomerView: fsCheck.Checked,
				ShowDeleteConfirmation: dnsaCheck.Checked,
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
