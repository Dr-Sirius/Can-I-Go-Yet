package converter

import (
	"can-i-go-yet/src/scheduler"
	"can-i-go-yet/src/templater"
	"image/color"

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

func ColorToInt(c color.Color) [4]int {
	r,g,b,a := c.RGBA()
	return [4]int{int(r),int(g),int(b),int(a)}
}

func IntToColor(i [4]int) color.Color {
	return color.RGBA{uint8(i[0]),uint8(i[1]),uint8(i[2]),uint8(i[3])}
}