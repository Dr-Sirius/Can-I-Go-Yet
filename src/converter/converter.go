package converter

import (
	"can-i-go-yet/src/schedules"
	"image/color"

	"fyne.io/fyne/v2/data/binding"
)

func DataItemToSchedule(d binding.DataItem) schedules.Schedule {
	s, _ := d.(binding.Untyped).Get()
	return s.(schedules.Schedule)
}


func ColorToInt(c color.Color) [4]int {
	r, g, b, a := c.RGBA()
	return [4]int{int(r), int(g), int(b), int(a)}
}

func IntToColor(i [4]int) color.Color {
	return color.RGBA{uint8(i[0]), uint8(i[1]), uint8(i[2]), uint8(i[3])}
}
