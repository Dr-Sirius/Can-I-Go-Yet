package templates

import (
	"can-i-go-yet/src/converter"
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/schedules"
	"can-i-go-yet/src/templates"
	"errors"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//region Build Templates UI

/*
Creates *widget.Form for creating new template
*/
func TemplateForm(list *widget.List, b *binding.UntypedList, win fyne.Window) *fyne.Container {
	// Template Name
	tName := widget.NewEntry()
	// Starting Time
	stEntry := widget.NewEntry()
	stEntry.SetText("12:00 am")
	// End Time
	etEntry := widget.NewEntry()
	etEntry.SetText("12:00 pm")

	saveBTN := widget.NewButtonWithIcon("Save Template", theme.DocumentSaveIcon(), func() {
		schedule := []schedules.Schedule{}
		for i := range (*b).Length() {
			item, _ := (*b).GetItem(i)
			schedule = append(schedule, converter.DataItemToSchedule(item))
		}
		if templates.Exists(tName.Text) {
			log.Println("true")
			dialog.NewError(errors.New("Template: "+tName.Text+" already exists"), win).Show()
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
			return
		}

		templates.CreateTemplateFile(templates.Template{Name: tName.Text, Schedules: schedule})
		log.Println("Created")
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
			(*b).Append(schedules.New(stEntry.Text, etEntry.Text, time.Now().Format("2006-01-02"), handler.CreateFlags(flags.Selected)))

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
