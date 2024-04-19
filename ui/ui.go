package ui

import (
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/templater"
	"can-i-go-yet/src/converter"

	"image/color"
	"time"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())
	
	myWindow := app.NewWindow("Can I Go Yet?")
	
	
	b := binding.NewUntypedList()
	for _,x := range scheduler.LoadSchedules(){
		b.Append(x)
	}

	content := container.NewAppTabs(
		container.NewTabItem("Today", TodayTab(b)),
		container.NewTabItem("Announcments", Announcments()),
		container.NewTabItem("Add Schedule", AddForm(b)),
		container.NewTabItem("Remove Schedule", Remove(b)),
		container.NewTabItem("Templates", TemplatTab()),
		container.NewTabItem("Build Template", BuildTemplatTab()),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(200, 200))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()

}

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

func AddForm(data binding.UntypedList) *widget.Form {
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
			s := scheduler.NewSchedule(stEntry.Text, etEntry.Text, dtEntry.Text,handler.CreateFlags(flags.Selected)...)
			scheduler.AddSchedule(dtEntry.Text, stEntry.Text, etEntry.Text, handler.CreateFlags(flags.Selected)...)
			data.Append(s)
			// if dtEntry.Text == handler.GetDate() {
			// 	handler.SetTime()
			// }
		},
	}
}

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
		s,_ := data.Get()
		temp := s[selected+1:]
		s = append(s[:selected], temp...)
		data.Set(s)
		selected = -1
		rl.UnselectAll()
		rl.Refresh()
	})

	removeAllBTN := widget.NewButton("Remove All", func() {
		handler.RemoveAll()
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

func TemplatTab() *fyne.Container {
	todayLBL := canvas.NewText("Templates", color.Black)
	todayLBL.TextSize = 35
	tabs := container.NewDocTabs(TemplateTabs()...)
	dateENT := widget.NewEntry()
	dateENT.SetText(time.Now().Format("2006-01-02"))

	name := tabs.Items[0].Text

	tabs.OnSelected = func(ti *container.TabItem) {
		name = ti.Text
	}

	tabs.OnClosed = func (ti *container.TabItem)  {
		templater.RemoveTemplate(ti.Text)
		
	}
	
	addBTN := widget.NewButton("Add Template", func() {

		for _, x := range templater.LoadTemplate(name) {

			scheduler.AddSchedule(dateENT.Text, x.Start_Time, x.End_Time, x.FlagsSlice()...)
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



func TemplateTabs() []*container.TabItem {
	var tabs []*container.TabItem
	
	for i, x := range templater.GetAllTemplates() {
		c := container.NewTabItem(
			i,
			TemplateList(x),
		)

		tabs = append(tabs, c)
	}
	
	return handler.SortTabs(tabs)
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

func TemplateForm(list *widget.List,b *binding.UntypedList) *widget.Form {
	tName := widget.NewEntry()
	tName.SetPlaceHolder("Alpha")
	stEntry := widget.NewEntry()
	stEntry.SetPlaceHolder("12:00 am")
	etEntry := widget.NewEntry()
	etEntry.SetPlaceHolder("12:00 pm")
	saveBTN := widget.NewButton("Save",func() {
		t := []templater.Template{}
		for i := range (*b).Length() {
			item,_ := (*b).GetItem(i)
			t = append(t, converter.DataItemToTemplate(item))
		}
		templater.AddTemplate(t)
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

func BuildTemplateList(data binding.UntypedList) *widget.List {
	return widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			lbl := canvas.NewText("template",color.Black)
			lbl.TextSize = 15	
			return lbl
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			
			o.(*canvas.Text).Text = converter.DataItemToTemplate(di).PrettyString()
		},
	)
}

func BuildTemplatTab() *container.Split {
	
	b := binding.NewUntypedList()

	ls := BuildTemplateList(b)

	content := container.NewVSplit(
		ls,	
		TemplateForm(ls,&b),
	)
	
	return content
}
