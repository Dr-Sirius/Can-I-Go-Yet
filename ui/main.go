package main

import (
	uic "can-i-go-yet/ui/customer"
	uiu "can-i-go-yet/ui/user"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	uic.CustomerView(myApp)
	uiu.UserView(myApp)
	
	myApp.Run()
}
