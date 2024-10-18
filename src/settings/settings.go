package settings

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type Settings struct {
	DefaultTemplate        string    `json:"DefaultTemplate"`
	StayOpen               bool      `json:"StayOpen"`
	OpenColor              [4]int    `json:"OpenColor"`
	ClosedColor            [4]int    `json:"ClosedColor"`
	BreakColor             [4]int    `json:"BreakColor"`
	StandardHours          [2]string `json:"StandarHours"`
	FullscreenCustomerView bool      `json:"FullscreenCustomerView"`
	ShowDeleteConfirmation bool      `json:"ShowDeleteConfirmation"`
}

func LoadSettings() Settings {
	js, err := os.Open("Settings/Settings.json")
	if err != nil {
		log.Println(err)
	}
	defer js.Close()
	s := Settings{}
	b, err := io.ReadAll(js)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(b, &s)
	return s
}

func (s Settings) SaveSettings() {
	ms, err := json.MarshalIndent(s, "", "	")
	if err != nil {
		log.Println(err)
	}
	if err = os.WriteFile("Settings/Settings.json", ms, os.ModePerm); err != nil {
		log.Println(err)
	}

}

func CreateSettings() {
	s := Settings{
		DefaultTemplate:        "TempTemplate",
		StayOpen:               true,
		OpenColor:              [4]int{0, 255, 0, 255},
		ClosedColor:            [4]int{255, 0, 0, 255},
		BreakColor:             [4]int{0, 0, 255, 255},
		StandardHours:          [2]string{"7:45 am", "2:30 pm"},
		FullscreenCustomerView: false,
		ShowDeleteConfirmation: true,
	}
	ms, err := json.MarshalIndent(s, "", "	")
	if err != nil {
		log.Println(err)
	}
	if _, err := os.Stat("Settings/Settings.json"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("Settings", os.ModePerm); err != nil {
			log.Println(err)
		} else {
			os.Create("Settings/Settings.json")
			os.WriteFile("Settings/Settings.json", ms, os.ModePerm)
		}
	}
}
