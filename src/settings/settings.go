package settings

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Settings struct {
	DefaultTemplate string    `json:"DefaultTemplate"`
	StayOpen        bool      `json:"StayOpen"`
	OpenColor       [4]int    `json:"OpenColor"`
	ClosedColor     [4]int    `json:"ClosedColor"`
	BreakColor      [4]int    `json:"BreakColor"`
	StandardHours   [2]string `json:"StandarHours"`
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
	json.Unmarshal(b,&s)
	return s
}

func (s Settings) SaveSettings() {
	ms, err := json.MarshalIndent(s,"","	")
	if err != nil {
		log.Println(err)
	}
	if err = os.WriteFile("../../Settings/Settings.json",ms,os.ModePerm); err != nil {
		log.Println(err)
	}

}
