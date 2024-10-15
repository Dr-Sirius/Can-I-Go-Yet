package templates

import (
	"can-i-go-yet/src/schedules"
	"log"
	"testing"
)

func TestCreateTemplateFile(t *testing.T) {
	s1 := schedules.New("12:00 am", "12:00 pm", "2024-08-19", []int{0})
	log.Println(s1.Flags)
	s2 := schedules.New("12:00 am", "12:00 pm", "2024-08-19", []int{0, 1, 2})
	temp := New("Alpha", []schedules.Schedule{s1, s2})
	CreateTemplateFile(temp)
}

func TestLoadTemplate(t *testing.T) {
	temp, _ := LoadTemplate("Alpha")
	log.Println(temp.Schedules[0].Date(), "before")
	temp.SetSchedulesForToday()
	log.Println(temp.Schedules[0].Date(), "after")
	log.Fatal()

}
