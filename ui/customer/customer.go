package ui

import (
	ui "can-i-go-yet/ui/user"
	_ "fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/eiannone/keyboard"
)


var updateTime time.Duration = 5 * time.Second

func CustomerView(app fyne.App) {

	myWindow := app.NewWindow("Customer View")
	myWindow.Resize(fyne.NewSize(800,600))

	openLBL := canvas.NewText(ui.UserInput,color.RGBA{R:255,G:0,B:0,A:255})
	openLBL.TextSize = 100
	openLBL.Alignment = fyne.TextAlignCenter
	content := container.New(layout.NewVBoxLayout(),openLBL)
	myWindow.SetContent(content)
	go update(content)
	myWindow.Show()

}

func update(content *fyne.Container) {

	for range time.Tick(updateTime) {
		
	}
}
