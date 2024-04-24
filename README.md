# Can I Go Yet
Can I Go Yet is an app designed to act as a display kiosk for displaying the staus and hours of a tech office

# Getting Started
To run it you may require specific packages if your running linux. The packages can be found at https://docs.fyne.io/started/ 

# Build
To build it your will need to make sure you have all the prerequisites for your platform. These can be found at https://docs.fyne.io/started/ 

To build it as a desktop app you will need to install fyne and run
```
fyne package --os {Your Operating System} --src app --icon ../Icon.png
```
If you`re using linux you will need to unzip the tar.gz and run the make file

if you don't want to use fyne you can simply run
```
go build app/main.go
```
and it will create the proper executable for your system
