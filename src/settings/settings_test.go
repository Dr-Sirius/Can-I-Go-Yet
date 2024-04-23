package settings

import (
	"log"
	"testing"
)

func TestSaveSettings(t *testing.T) {
	s := Settings{
		DefaultTemplate: "Demo",
		StayOpen:        true,
		OpenColor:       [4]int{0, 255, 0, 255},
		ClosedColor:     [4]int{255, 0, 0, 255},
		BreakColor:      [4]int{0, 0, 255, 255},
		StandardHours:   [2]string{"7:00 am", "2:30 pm"},
	}
	s.SaveSettings()
}

func TestLoadSettings(t *testing.T) {

	log.Println(LoadSettings())
}
