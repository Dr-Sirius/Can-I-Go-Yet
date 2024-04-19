package converter

import (
	"can-i-go-yet/src/templater"

	"fyne.io/fyne/v2/data/binding"
)

func DataItemToTemplate(d binding.DataItem) templater.Template {
	s,_ := d.(binding.Untyped).Get()
	return s.(templater.Template)
}