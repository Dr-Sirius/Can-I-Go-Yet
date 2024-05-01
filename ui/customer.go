package ui

import (
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/settings"
	"fmt"

	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var previousAnouncment string = handler.Announcments

func CustomerView() {
	set := settings.LoadSettings()
	myWindow := fyne.CurrentApp().NewWindow("Customer View")
	//myWindow.SetFullScreen(set.FullscreenCustomerView)

	openLBL := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	openLBL.TextSize = 150
	openLBL.Alignment = fyne.TextAlignCenter

	ctLBL := canvas.NewText("", color.Black)
	ctLBL.TextSize = 50
	ctLBL.Alignment = fyne.TextAlignCenter

	dHours := fmt.Sprintf("Tech Office Daily Hours: %s - %s", set.StandardHours[0], set.StandardHours[1])
	officeHoursLBL := canvas.NewText(dHours, color.Black)
	officeHoursLBL.TextSize = 50
	officeHoursLBL.TextStyle.Bold = true
	officeHoursLBL.Alignment = fyne.TextAlignCenter

	statusLBL := canvas.NewText("", color.Black)
	statusLBL.TextSize = 50
	statusLBL.Alignment = fyne.TextAlignCenter

	// announcmentsLBL := canvas.NewText("Announcments:", color.Black)
	// announcmentsLBL.TextSize = 50
	// announcmentsLBL.Alignment = fyne.TextAlignCenter

	announcmentsBODY := canvas.NewText("", color.Black)
	announcmentsBODY.TextSize = 25
	announcmentsBODY.TextStyle.Bold = true
	announcmentsBODY.Alignment = fyne.TextAlignCenter

	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillOriginal

	content := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			logo,
			layout.NewSpacer(),
			openLBL,
			layout.NewSpacer(),
		),

		widget.NewSeparator(),

		container.NewHBox(
			ctLBL,
			widget.NewSeparator(),
			announcmentsBODY,
		),

		widget.NewSeparator(),

		statusLBL,
	)

	myWindow.SetContent(content)
	go func() {
		updateClock(ctLBL)
		for range time.Tick(time.Second) {
			updateClock(ctLBL)
			updateOpen(openLBL)
			updateStatus(statusLBL)
			updateAnnouncments(announcmentsBODY)
			set = settings.LoadSettings()
			dHours := fmt.Sprintf("Tech Office Daily Hours: %s - %s", set.StandardHours[0], set.StandardHours[1])
			officeHoursLBL.Text = dHours
			content.Refresh()
		}
	}()
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.Show()

}

func updateClock(clock *canvas.Text) {
	currentTime := time.Now().Format("03:04:05 pm")
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
		status.Text = "Will reopen at " + handler.GetReturnTime()
		return
	}
	if handler.Status == "Break" {
		status.Text = "There may be increased wait times"
		return
	}
	status.Text = ""
}

func updateAnnouncments(anc *canvas.Text) {
	if previousAnouncment != handler.Announcments {
		anc.TextSize -= float32(len(anc.Text))
		previousAnouncment = handler.Announcments
	}
	anc.Text = handler.Announcments
	if len(anc.Text) == 0 {
		anc.TextSize = 50
	}
}
