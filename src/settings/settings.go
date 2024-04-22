package settings

type Settings struct {
	DefaultTemplate string `json:"DefaultTemplate"`
	StayOpen        bool   `json:"StayOpen"`
	OpenColor       [4]int `json:"OpenColor"`
	ClosedColor     [4]int `json:"OpenColor"`
	BreakColor      [4]int
	StandardHours   [2]string
}

func LoadSettings() Settings

func (s Settings) SaveSettings() {

}
