package ui

import (
	"can-i-go-yet/src/scheduler"
	_ "fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
)

var sch []scheduler.Schedule
var currentSch int
var stayOpen bool = false

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
	officeHoursLBL.Alignment = fyne.TextAlignCenter

	logo := canvas.NewImageFromResource(resourceLogoPng)

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
	go func() {
		updateClock(ctLBL)
		setTime()
		for range time.Tick(time.Second) {
			updateClock(ctLBL)
			checkTime(openLBL)
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

func setTime() {
	for i, x := range sch {
		if time.Now().Equal(x.StartTime) {
			currentSch = i
			break
		}
	}
}

func checkTime(txt *canvas.Text) {
	if len(sch) == 0 {
		log.Println("SCH = 0")
		setOpen(txt)

	} else if time.Now().Equal(sch[currentSch].StartTime) || time.Now().After(sch[currentSch].StartTime) {
		log.Println("Starttime")
		checkFlags(txt)

	} else if time.Now().Equal(sch[currentSch].EndTime) || time.Now().After(sch[currentSch].EndTime) {
		log.Println("Endtime")
		if (currentSch + 1) != len(sch) {

			currentSch += 1
			checkTime(txt)
		} else {
			setClosed(txt)

		}

	} else {
		if stayOpen {
			setOpen(txt)
		} else {
			setClosed(txt)
		}
	}

	txt.Refresh()
}

func checkFlags(txt *canvas.Text) {
	log.Println(txt)
	flags := sch[currentSch].Flags
	log.Println(flags)
	if _, ok := flags[scheduler.BRKE]; ok {
		if _, ok := flags[scheduler.UNDS]; ok {
			setClosed(txt)
		} else {
			setOnBreak(txt)
		}

	} else if _, ok := flags[scheduler.OPEN]; ok {
		setOpen(txt)
	} else {
		setClosed(txt)
	}

}

func setOpen(txt *canvas.Text) {
	txt.Text = "Open"
	txt.Color = color.RGBA{0,255,0, 255}
}

func setClosed(txt *canvas.Text) {
	txt.Text = "Closed"
	txt.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
}

func setOnBreak(txt *canvas.Text) {
	txt.Text = "Open"
	txt.Color = color.RGBA{255, 255, 0, 255}
}
