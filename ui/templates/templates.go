package templates

import (
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/schedules"
	"can-i-go-yet/src/settings"
	"can-i-go-yet/src/templates"
	"errors"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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

	})

	replaceBTN := widget.NewButton("Replace Current Schedules with Template", func() {
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
		handler.RemoveSchedulesFromCurrent()
		
		template, _ := templates.LoadTemplate(name)
		
		template.SetSchedulesForToday()

		handler.AssignNewSchedules(template)
		for _, x := range template.Schedules {

			data.Append(x)
			
		}

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
