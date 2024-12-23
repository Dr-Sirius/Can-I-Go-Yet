package main

import (
	"can-i-go-yet/src"
	"can-i-go-yet/ui"
)

func main() {

	go src.Start()

	ui.Run()

}
