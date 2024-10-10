package templates

import (
	"can-i-go-yet/src/schedules"
	"encoding/json"
	"errors"

	"log"
	"os"
	
)

/*
Struct for holding template information
*/
type Template struct {
	Name       string `json:"Name"`
	Schedules 	[]schedules.Schedule `json:"Schedules"`
}

/*
Creates and returns new template based on params
*/
func New(name string, schedules []schedules.Schedule) Template {
	return Template{Name: name, Schedules: schedules}
}


func CreateTemplateFile(template Template) {
	
	templateJson, err := json.MarshalIndent(template, "", "	")
	if err != nil {
		log.Println(err)
	}
	if _, err := os.Stat("Templates/t_" + template.Name + ".json"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("Templates", os.ModePerm); err != nil {
			log.Println(err)
		} else {
			os.Create("Templates/t_" + template.Name + ".json")
			os.WriteFile("Templates/t_" + template.Name + ".json", templateJson, os.ModePerm)
		}
	}
}

/*
Creates Templates folder
*/
func CreateTemplatesFolder() {
	if _, err := os.Stat("Templates"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("Templates", os.ModePerm); err != nil {
			log.Println(err)
		} else {
			os.Create("Templates")
		}
	}
	
}


/*
Removes specified template
*/
func RemoveTemplate(name string) error {
	return os.Remove("Templates/t_" + name + ".json")
}


/*
Returns boolean based on if passed template exists
*/
// func Exists(name string) bool {
// 	_,ok := GetAllTemplates()[name]

// 	return ok
// }
