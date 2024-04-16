package templater

import (
	"log"
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	//log.Println(CreateTemplate("beta"))
	tp := []Template{
		NewTemplate("Death", "7:00 am", "10:30 am", 0),
		NewTemplate("Death", "11:00 am", "11:30 am", 0),
		NewTemplate("Death", "12:00 pm", "12:30 pm", 0),
		NewTemplate("Death", "1:00 pm", "1:30 pm", 0),
	}
	AddTemplate(tp)
	for _, x := range LoadTemplate(tp[0].Name) {
		log.Println(x)
	}
}

func TestGetAllTemplates(t *testing.T) {
	log.Println(GetAllTemplates())
}
