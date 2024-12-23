package ui

import (
	"can-i-go-yet/src/handler"
	"can-i-go-yet/ui/schedules"
	"can-i-go-yet/ui/settings"
	"can-i-go-yet/ui/templates"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
)

/*
Creates and Runs a new fyne.App - handles all ui related events
*/
func Run() {
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())

	myWindow := app.NewWindow("Can I Go Yet?")

	b := binding.NewUntypedList()
	for _, x := range handler.GetCurrentSchedules() {
		b.Append(x)
	}

	content := container.NewAppTabs(
		container.NewTabItem("Today", schedules.TodayTab(b)),
		container.NewTabItem("Announcments", schedules.Announcments()),
		container.NewTabItem("Remove Schedule", schedules.Remove(b, myWindow)),
		container.NewTabItem("Templates", templates.TemplateTab(b, myWindow)),
		container.NewTabItem("Build Template", templates.BuildTemplatTab(myWindow)),
		container.NewTabItem("Settings", settings.SettingsTab(myWindow)),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()

}
