package ui

import (
	
	_ "fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//var updateTime time.Duration = 5 * time.Second

func CustomerView(app fyne.App) {

	myWindow := app.NewWindow("Customer View")
	myWindow.Resize(fyne.NewSize(800, 600))
	

	openLBL := canvas.NewText("Closed", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	openLBL.TextSize = 150
	openLBL.Alignment = fyne.TextAlignCenter

	ctLBL := canvas.NewText("", color.Black)
	ctLBL.TextSize = 25
	ctLBL.Alignment = fyne.TextAlignCenter

	officeHoursLBL := canvas.NewText("Tech Office Daily Hours: 7:40 am - 2:40 pm", color.Black)
	officeHoursLBL.TextSize = 25
	officeHoursLBL.Alignment = fyne.TextAlignCenter

	logo := canvas.NewImageFromResource(resourceLogoPng)

	

	go func() {
		updateClock(ctLBL)
		for range time.Tick(time.Second) {
			updateClock(ctLBL)
		}
	}()

	content := container.New(
		layout.NewAdaptiveGridLayout(3),
		logo,
		container.New(
			layout.NewVBoxLayout(),
			openLBL,
			ctLBL,
			officeHoursLBL,
		),
	)
	myWindow.SetContent(content)
	myWindow.Show()

}

func updateClock(clock *canvas.Text) {
	currentTime := time.Now().Format("Current Time: 03:04:05 AM")
	clock.Text = currentTime
	clock.Refresh()

}

func NewFullScreenButton(text string) *widget.Label {
	return widget.NewLabel(text)
}
