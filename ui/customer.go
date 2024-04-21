package ui

import (
	"can-i-go-yet/src/handler"

	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/layout"
)

func CustomerView() {

	myWindow := fyne.CurrentApp().NewWindow("Customer View")
	openLBL := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	openLBL.TextSize = 150
	openLBL.Alignment = fyne.TextAlignCenter

	ctLBL := canvas.NewText("", color.Black)
	ctLBL.TextSize = 25
	ctLBL.Alignment = fyne.TextAlignCenter

	officeHoursLBL := canvas.NewText("Tech Office Daily Hours: 7:40 am - 2:40 pm", color.Black)
	officeHoursLBL.TextSize = 25
	officeHoursLBL.TextStyle.Bold = true
	officeHoursLBL.Alignment = fyne.TextAlignCenter

	statusLBL := canvas.NewText("",color.Black)
	statusLBL.TextSize = 25
	statusLBL.Alignment = fyne.TextAlignCenter

	announcmentsLBL := canvas.NewText("Announcments:",color.Black)
	announcmentsLBL.TextSize = 25
	announcmentsLBL.Alignment = fyne.TextAlignCenter

	announcmentsBODY := widget.NewMultiLineEntry()
	announcmentsBODY.Text = "Bacon"
	announcmentsBODY.TextStyle.Bold = true

	logo := canvas.NewImageFromResource(resourceLogoPng)

	content := container.New(
		layout.NewAdaptiveGridLayout(2),
		logo,
		container.New(
			layout.NewVBoxLayout(),
			openLBL,
			canvas.NewLine(color.Black),
			ctLBL,
			officeHoursLBL,
			canvas.NewLine(color.Black),
			statusLBL,
			announcmentsLBL,
			announcmentsBODY,
		),
	)

	myWindow.SetContent(content)
	go func() {
		updateClock(ctLBL)
		for range time.Tick(time.Second) {
			updateClock(ctLBL)
			updateOpen(openLBL)
			updateStatus(statusLBL)
			updateAnnouncments(announcmentsBODY)
			content.Refresh()
		}
	}()
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.Show()

}

func updateClock(clock *canvas.Text) {
	currentTime := time.Now().Format("Current Time: 03:04:05 pm")
	clock.Text = currentTime
	clock.Refresh()
}

func updateOpen(open *canvas.Text) {
	status, colour := handler.CheckTime()
	open.Text = status
	open.Color = colour
}

func updateStatus(status *canvas.Text) {
	if handler.Status == "Closed" {
		status.Text = "The Tech Office will reopen at " + handler.GetReturnTime()
		return
	}
	if handler.Status == "Break" {
		status.Text = "There may be increased wait times"
		return
	}
	status.Text = ""
}

func updateAnnouncments(anc *widget.Entry) {
	anc.Text = handler.Announcments
}