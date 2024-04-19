package templater

import (
	"can-i-go-yet/src/scheduler"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Template struct {
	Name       string
	Start_Time string
	End_Time   string
	Flags      map[int]bool
}


func NewTemplate(name string, startTime string, endTime string, flags ...int) Template {
	f := make(map[int]bool)
	for _, x := range flags {
		f[x] = true
	}
	return Template{Name: name, Start_Time: startTime, End_Time: endTime, Flags: f}
}



func CreateTemplate(Name string) error {
	if _, err := os.Stat("Templates"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("Templates", os.ModePerm); err != nil {
			return err
		} else {
			os.WriteFile("Templates/t_"+Name+".csv", []byte("Name, Start_Time, End_Time, Flags"), os.ModePerm)
			return nil
		}
	}
	return nil
}

func AddTemplate(t []Template) {
	os.WriteFile("Templates/t_"+t[0].Name+".csv", []byte("Name, Start_Time, End_Time, Flags"), os.ModePerm)
	for _, x := range t {
		flagString := ""
		for _, f := range x.FlagsSlice() {
			flagString += fmt.Sprint(f) + "|"
		}

		sch := "\n" + x.Name + ", " + x.Start_Time + ", " + x.End_Time + ", " + flagString

		file, err := os.OpenFile("Templates/t_"+x.Name+".csv", os.O_APPEND|os.O_WRONLY, 0644)

		if err != nil {
			log.Println(err)
		}

		defer file.Close()

		_, err = file.WriteString(sch)

		if err != nil {
			log.Println(err)

		}
	}

}

func AddTemplateWithName(t []Template, name string) {
	os.WriteFile("Templates/t_"+name+".csv", []byte("Name, Start_Time, End_Time, Flags"), os.ModePerm)
	for _, x := range t {
		flagString := ""
		for _, f := range x.FlagsSlice() {
			flagString += fmt.Sprint(f) + "|"
		}

		sch := "\n" + x.Name + ", " + x.Start_Time + ", " + x.End_Time + ", " + flagString

		file, err := os.OpenFile("Templates/t_"+name+".csv", os.O_APPEND|os.O_WRONLY, 0644)

		if err != nil {
			log.Println(err)
		}

		defer file.Close()

		_, err = file.WriteString(sch)

		if err != nil {
			log.Println(err)

		}
	}
}

func LoadTemplate(name string) []Template {
	if err := CreateTemplate(name); err != nil {
		log.Println(err)
	}
	file, err := os.Open("Templates/t_" + name + ".csv")

	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		log.Println(err)
	}

	var t []Template

	for _, r := range records[1:] {
		f := strings.Split(r[3], "|")
		var fs []int
		for _, x := range f {
			cf, err := strconv.Atoi(strings.Trim(x, " "))
			if err != nil {
				log.Println(err)
			}
			fs = append(fs, cf)
		}
		ns := NewTemplate(r[0], r[1], r[2], fs...)
		t = append(t, ns)
	}
	return t
}

func RemoveTemplate(name string) error {
	return os.Remove("Templates/t_" + name + ".csv")
}

func GetAllTemplates() map[string][]Template {
	entries, err := os.ReadDir("Templates")
	if err != nil {
		log.Fatal(err)
	}
	t := make(map[string][]Template)

	for _, e := range entries {
		log.Println(e.Name())
		tp := LoadTemplate(e.Name()[2 : len(e.Name())-4])
		t[e.Name()[2:len(e.Name())-4]] = tp

	}
	return t
}

func (t Template) PrettyString() string {
	st := t.Start_Time
	et := t.End_Time
	f := func() string {
		str := ""
		for i := range t.Flags {
			if i == scheduler.OPEN {
				str += "|Open"
			}
			if i == scheduler.BRKE {
				str += "|Break"
			}
			if i == scheduler.UNDS {
				str += "|Understaffed"
			}
			if i == scheduler.HDAY {
				str += "|Holiday"
			}
		}
		return str
	}
	return st + " - " + et + " " + f()
}

func (s Template) FlagsSlice() []int {
	x := []int{}
	for a, _ := range s.Flags {
		x = append(x, a)
	}
	slices.Sort(x)
	return x
}

func Exists(name string) bool {
	_,ok := GetAllTemplates()[name]

	return ok
}
