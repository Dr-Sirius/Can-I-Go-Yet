package main

import (
	"can-i-go-yet/ui"
	
)


//go:generate fyne bundle -o ui/logo.go -append ui/assets/logo.png

func main() {

	ui.Run()
}
