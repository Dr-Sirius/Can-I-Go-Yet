package ui

import (
	_"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_"fyne.io/fyne/v2/theme"
)




func Run() {
	myApp := app.New()
	myApp.Settings().SetTheme(&defaultTheme{})
	UserView(myApp)
	myApp.Run()
}
