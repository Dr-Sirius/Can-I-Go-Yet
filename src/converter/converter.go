package converter

import (
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/templater"

	"fyne.io/fyne/v2/data/binding"
)

func DataItemToTemplate(d binding.DataItem) templater.Template {
	s, _ := d.(binding.Untyped).Get()
	return s.(templater.Template)
}

func DataItemToSchedule(d binding.DataItem) scheduler.Schedule {
	s, _ := d.(binding.Untyped).Get()
	return s.(scheduler.Schedule)
}

func TemplateToSchedule(name string, date string) []scheduler.Schedule {
	s := []scheduler.Schedule{}

	for _, x := range templater.LoadTemplate(name) {
		s = append(s, scheduler.NewSchedule(x.Start_Time, x.End_Time, date, x.FlagsSlice()...))
	}
	return s
}
