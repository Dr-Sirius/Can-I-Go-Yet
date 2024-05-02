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
	
	myWindow := fyne.CurrentApp().NewWindow("Customer View")
	myWindow.SetFullScreen(settings.LoadSettings().FullscreenCustomerView)


	
	go func() {
		
		for range time.Tick(time.Second) {
			myWindow.SetFullScreen(settings.LoadSettings().FullscreenCustomerView)
			myWindow.SetContent(CustomerContent())
			myWindow.Content().Refresh()
			
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


func CustomerContent() *fyne.Container {

	set := settings.LoadSettings()
	openLBL := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	openLBL.TextSize = 150
	openLBL.Alignment = fyne.TextAlignCenter
	openLBL.TextStyle.Bold = true

	ctLBL := canvas.NewText("", color.Black)
	ctLBL.TextSize = 70
	ctLBL.Alignment = fyne.TextAlignCenter

	dHours := fmt.Sprintf("Tech Office Daily Hours: %s - %s", set.StandardHours[0], set.StandardHours[1])
	officeHoursLBL := canvas.NewText(dHours, color.Black)
	officeHoursLBL.TextSize = 50
	officeHoursLBL.TextStyle.Bold = true
	officeHoursLBL.Alignment = fyne.TextAlignCenter

	statusLBL := canvas.NewText("", color.Black)
	statusLBL.TextSize = 70
	statusLBL.Alignment = fyne.TextAlignCenter

	// announcmentsLBL := canvas.NewText("Announcments:", color.Black)
	// announcmentsLBL.TextSize = 50
	// announcmentsLBL.Alignment = fyne.TextAlignCenter

	announcmentsBODY := canvas.NewText("", color.Black)
	announcmentsBODY.TextSize = 25
	announcmentsBODY.TextStyle.Bold = true
	//announcmentsBODY.Alignment = fyne.TextAlignCenter

	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillOriginal


	updateAnnouncments(announcmentsBODY)
	updateClock(ctLBL)
	updateOpen(openLBL)
	updateStatus(statusLBL)

	if announcmentsBODY.Text == "" {
		ctLBL.TextSize = 80
		statusLBL.TextSize = 80
		return container.NewVBox(
			container.NewHBox(
				layout.NewSpacer(),
				logo,
				layout.NewSpacer(),
				openLBL,
				layout.NewSpacer(),
			),
	
			widget.NewSeparator(),
	
			container.NewHBox(
				layout.NewSpacer(),
				ctLBL,
				layout.NewSpacer(),
			),
			
			widget.NewSeparator(),
	
			statusLBL,
		)
	}

	return container.NewVBox(
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
			layout.NewSpacer(),
			announcmentsBODY,
			layout.NewSpacer(),
		),
		widget.NewSeparator(),

		statusLBL,
	)
}