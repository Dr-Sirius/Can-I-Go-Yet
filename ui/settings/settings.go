package settings

import (
	"can-i-go-yet/src/converter"
	"can-i-go-yet/src/handler"
	"can-i-go-yet/src/settings"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//region Settings UI

/*
Creates *widget.Form for displaying and changing settings - reads and writes to Settings/Settings.json
*/
func SettingsTab(w fyne.Window) *widget.Form {
	tName := widget.NewEntry()
	tName.Text = handler.GetDefaultTemplate()
	stoCheck := widget.NewCheck("", func(b bool) { handler.SetStayOpen(b) })

	var oColor color.Color = converter.IntToColor(settings.LoadSettings().OpenColor)
	opRect := canvas.NewRectangle(oColor)
	var bColor color.Color = converter.IntToColor(settings.LoadSettings().BreakColor)
	bkRect := canvas.NewRectangle(bColor)
	var cColor color.Color = converter.IntToColor(settings.LoadSettings().ClosedColor)
	clRect := canvas.NewRectangle(cColor)

	openColorDialog := dialog.NewColorPicker("Open Color", "", func(c color.Color) { oColor = c; SetColor(c, opRect) }, w)
	openColorDialog.Advanced = true
	closedColorDialog := dialog.NewColorPicker("Closed Color", "", func(c color.Color) { cColor = c; SetColor(c, clRect) }, w)
	closedColorDialog.Advanced = true
	breakColorDialog := dialog.NewColorPicker("Break Color", "", func(c color.Color) { bColor = c; SetColor(c, bkRect) }, w)
	breakColorDialog.Advanced = true

	opBTN := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), openColorDialog.Show)
	clBTN := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), closedColorDialog.Show)
	bkBTN := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), breakColorDialog.Show)

	oContent := container.NewGridWithColumns(2, opRect, opBTN)
	cContent := container.NewGridWithColumns(2, clRect, clBTN)
	bContent := container.NewGridWithColumns(2, bkRect, bkBTN)

	stEntry := widget.NewEntry()
	stEntry.SetText(settings.LoadSettings().StandardHours[0])

	etEntry := widget.NewEntry()
	etEntry.SetText(settings.LoadSettings().StandardHours[1])

	dhContent := container.NewVBox(stEntry, etEntry)

	fsCheck := widget.NewCheck("", func(b bool) {})
	fsCheck.SetChecked(settings.LoadSettings().FullscreenCustomerView)

	dnsaCheck := widget.NewCheck("", func(b bool) {})
	dnsaCheck.SetChecked(settings.LoadSettings().ShowDeleteConfirmation)

	return &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Default Template", Widget: tName},
			{Text: "Stay Open", Widget: stoCheck},
			{Text: "Open Label Color", Widget: oContent},
			{Text: "Closed Label Color", Widget: cContent},
			{Text: "Break Label Color", Widget: bContent},
			{Text: "Daily Hours", Widget: dhContent},
			{Text: "Fullscreen Customer View", Widget: fsCheck},
			{Text: "Show delete confirmation", Widget: dnsaCheck},
		},
		SubmitText: "Save",
		OnSubmit: func() {

			s := settings.Settings{
				DefaultTemplate:        tName.Text,
				StayOpen:               stoCheck.Checked,
				OpenColor:              converter.ColorToInt(oColor),
				ClosedColor:            converter.ColorToInt(cColor),
				BreakColor:             converter.ColorToInt(bColor),
				StandardHours:          [2]string{stEntry.Text, etEntry.Text},
				FullscreenCustomerView: fsCheck.Checked,
				ShowDeleteConfirmation: dnsaCheck.Checked,
			}
			s.SaveSettings()
			handler.Update()

		},
	}
}

/*
Sets passed *canvas.Rectange color to passed color.Color
*/
func SetColor(c color.Color, r *canvas.Rectangle) {
	r.FillColor = c
}
