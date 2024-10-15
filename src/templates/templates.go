package templates

import (
	"can-i-go-yet/src/schedules"
	"encoding/json"
	"errors"

	"io"
	"log"
	"os"
)

/*
Struct for holding template information
*/
type Template struct {
	Name      string               `json:"Name"`
	Schedules []schedules.Schedule `json:"Schedules"`
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
	if Exists("Templates/t_" + template.Name + ".json") {
		if err := os.Mkdir("Templates", os.ModePerm); err != nil {
			log.Println(err)
		} else {
			os.Create("Templates/t_" + template.Name + ".json")
			os.WriteFile("Templates/t_"+template.Name+".json", templateJson, os.ModePerm)
		}
	}
}

/*
Creates Templates folder
*/
func CreateTemplatesFolder() {
	if Exists("Templates") {
		if err := os.Mkdir("Templates", os.ModePerm); err != nil {
			log.Println(err)
		} else {
			os.Create("Templates")
		}
	}

}

/*
Returns the specified template from the templates folder
*/
func LoadTemplate(name string) (Template, error) {
	js, err := os.Open("Templates/t_" + name + ".json")
	if err != nil {
		return Template{}, err
	}
	defer js.Close()
	t := Template{}
	b, err := io.ReadAll(js)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(b, &t)
	return t, nil
}

/*
Removes specified template from Templates folder
*/
func RemoveTemplate(name string) error {
	return os.Remove("Templates/t_" + name + ".json")
}

/*
Changes the template's schedules with a version of the schedules with the current date
*/
func (t *Template) SetSchedulesForToday() {
	temp := make([]schedules.Schedule, len(t.Schedules))
	for i, s := range t.Schedules {
		temp[i] = schedules.GetScheduleWithTodayDate(s)
	}
	log.Println(temp)
	t.Schedules = temp
}

/*
Returns boolean based on if passed file exists
*/
func Exists(name string) bool {
	_, err := os.Stat("Templates")
	return errors.Is(err, os.ErrNotExist)
}
