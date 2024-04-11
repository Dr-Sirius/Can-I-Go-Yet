package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type defaultTheme struct{}

var _ fyne.Theme = (*defaultTheme)(nil)

func (t defaultTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	
	return color.White
}

func (t defaultTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t defaultTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (t defaultTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	
	return theme.DefaultTheme().Icon(name)
}




