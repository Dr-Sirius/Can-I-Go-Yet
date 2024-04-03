package ui

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	_"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	
)


var UserInput string = ""

func UserView(app fyne.App) {

	myWindow := app.NewWindow("Customer View")
	myWindow.Resize(fyne.NewSize(800,600))
	entry := widget.NewEntry()
	btn := widget.NewButton("enter",func() {
		UserInput = entry.Text
	})
	content := container.New(layout.NewVBoxLayout(),entry,btn)
	myWindow.SetContent(content)
	myWindow.Show()

}
