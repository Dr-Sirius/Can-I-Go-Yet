package schedules

import (
	"can-i-go-yet/src/converter"
	"can-i-go-yet/src/handler"
	"can-i-go-yet/ui/customer"

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

//region Schedules UI

/*
Creates a *widget.List that contains todays schedules
*/
func TodayList(data binding.UntypedList) *widget.List {

	return widget.NewListWithData(
		data,
		func() fyne.CanvasObject {

			lbl := canvas.NewText("template", color.Black)
			lbl.TextSize = 15
			return lbl
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = converter.DataItemToSchedule(di).PrettyString()
		},
	)
}

/*
Creates a *fyne.Container tab that contains information about the todays date including customer view
*/
func TodayTab(data binding.UntypedList) *fyne.Container {
	todayLBL := canvas.NewText("Today's Schedule", color.Black)
	todayLBL.TextSize = 35
	
	currentLBL := canvas.NewText("Current Schedule: "+handler.GetCurrentScheduleString(), color.Black)

	customerBTN := widget.NewButton("Customer View", func() {
		customer.CustomerView()
	})

	go func() {
		for range time.Tick(time.Second) {
			currentLBL.Text = "Current Schedule: " + handler.GetCurrentScheduleString()
		}
	}()
	return container.NewGridWithRows(
		4,
		todayLBL,
		TodayList(data),
		currentLBL,
		customerBTN,
	)

}

// /*
// Creates a *widget.Form for creating new schedule
// */
// func AddForm(data binding.UntypedList, win fyne.Window) *widget.Form {
// 	dtEntry := widget.NewEntry()
// 	dtEntry.SetText(time.Now().Format("2006-01-02"))
// 	stEntry := widget.NewEntry()
// 	stEntry.SetText("12:00 am")
// 	etEntry := widget.NewEntry()
// 	etEntry.SetText("12:00 pm")
// 	flags := widget.CheckGroup{
// 		Horizontal: true,
// 		Options: []string{
// 			"Open",
// 			"Break",
// 			"Understaffed",
// 			"Holiday",
// 			"Closed",
// 		},
// 	}
// 	return &widget.Form{
// 		Items: []*widget.FormItem{ // we can specify items in the constructor
// 			{Text: "Date", Widget: dtEntry},
// 			{Text: "Start Time", Widget: stEntry},
// 			{Text: "End Time", Widget: etEntry},
// 			{Text: "Flags", Widget: &flags},
// 		},
// 		OnSubmit: func() {
// 			if dtEntry.Text == "" || stEntry.Text == "" || etEntry.Text == "" {

// 				dialog.NewError(errors.New("the entries cannot be left blank"), win).Show()

// 				return

// 			}
// 			s := schedules.New(stEntry.Text, etEntry.Text, dtEntry.Text, handler.CreateFlags(flags.Selected))
// 			//schedules.AddSchedule(dtEntry.Text, stEntry.Text, etEntry.Text, handler.CreateFlags(flags.Selected)...)
// 			data.Append(s)
// 			if dtEntry.Text == handler.GetDate() {
// 				handler.SetTime()
// 			}
// 		},
// 	}
// }

/*
Creates *fyne.Container for removing schedules on todays date
*/
func Remove(data binding.UntypedList, win fyne.Window) *fyne.Container {
	selected := -1
	lbl := canvas.NewText("", color.Black)
	rl := TodayList(data)
	removeBTN := widget.NewButton("Remove", func() {
		if selected == -1 {

			dialog.NewError(errors.New("there are no schedules selected"), win).Show()

			return

		}

		handler.RemoveScheduleFromCurrent(selected)
		s, _ := data.Get()
		temp := s[selected+1:]
		s = append(s[:selected], temp...)
		data.Set(s)
		selected = -1
		rl.UnselectAll()
		rl.Refresh()
	})

	// removeAllBTN := widget.NewButton("Remove All", func() {
	// 	if len(schedules.LoadSchedules()) == 0 {

	// 		dialog.NewError(errors.New("there are no schedules\n\n you can create new schedules in the Add Schedule tab"),win).Show()

	// 		return

	// 	}
	// 	handler.RemoveAll()
	// 	s := make([]interface{}, 0)
	// 	data.Set(s)
	// 	rl.UnselectAll()
	// 	rl.Refresh()
	// })

	rl.OnSelected = func(id widget.ListItemID) {
		lbl.Text = "Remove " + handler.GetSchedule(id).PrettyString() + " ?"
		selected = id

		lbl.Refresh()
	}

	return container.NewGridWithRows(
		3,
		rl,
		lbl,
		container.NewHBox(
			removeBTN,
			//removeAllBTN,
		),
	)

}

/*
Creates *widget.Form for adding anouncment to customer view
*/
func Announcments() *widget.Form {
	anc := widget.NewEntry()

	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Anouncments", Widget: anc},
		},
		OnSubmit: func() {
			handler.Announcments = anc.Text
		},
	}
}
