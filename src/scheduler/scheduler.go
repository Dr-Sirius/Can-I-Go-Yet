package scheduler

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"strconv"
)

const (
	OPEN int = iota
	BRKE
	UNDS // understaffed
	HDAY // holiday
)

type Schedule struct {
	StartTime time.Time
	EndTime   time.Time
	Flags     map[int]bool
}

// Holds a slice of schedules for the specified date
type Day struct {
	// formated as YYYY-MM-DD
	Date      string
	Schedules []Schedule
}

func NewSchedule(st string, et string, date string, flags ...int) Schedule {
	return NewScheduleFromTime(convertTime(st, date), convertTime(et, date), flags...)
}

func NewScheduleFromTime(st time.Time, et time.Time, flags ...int) Schedule {
	f := make(map[int]bool)
	for _, x := range flags {
		f[x] = true
	}
	return Schedule{st, et, f}
}

func convertTime(s string, d string) time.Time {
	format := "2006-01-02 3:04 pm"
	s = d + " " + s
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println(err)
	}
	t, err := time.ParseInLocation(format, s, loc)
	if err != nil {
		log.Println(err)
	}
	return t
}

func NewDay(date string, schedules ...Schedule) Day {
	s := make([]Schedule, 0)
	for _, x := range schedules {
		y, m, d := x.StartTime.Date()
		mt := fmt.Sprint(int(m))
		if int(m) < 10 {
			mt = "0" + mt
		}

		dt := fmt.Sprint(y) + "-" + mt + "-" + fmt.Sprint(d)
		log.Println(dt)
		if dt == date {
			s = append(s, x)
		}
	}

	return Day{date, scheduleSort(s...)}
}

func scheduleSort(Schedules ...Schedule) []Schedule {
	for i := len(Schedules) - 1; i >= 0; i -= 1 {
		//log.Println(i)
		for x := range i {
			//log.Println(x)
			if Schedules[x].StartTime.Compare(Schedules[x+1].StartTime) > 0 {
				temp := Schedules[x]
				Schedules[x] = Schedules[x+1]
				Schedules[x+1] = temp
			}

		}
	}
	return Schedules
}

func (s Schedule) StringStartTime() string {
	h, m, _ := s.StartTime.Clock()
	tm := "am"

	if h > 12 {
		h = h - 12
		tm = "pm"
	}
	if h == 12 {
		tm = "pm"
	}

	hs := fmt.Sprint(h)
	ms := fmt.Sprint(m)

	if m < 10 && m != 0 {
		ms = "0" + ms
	}
	return hs + ":" + ms + " " + tm
}

func (s Schedule) StringEndTime() string {
	h, m, _ := s.EndTime.Clock()
	tm := "am"

	if h > 12 {
		h = h - 12
		tm = "pm"
	}
	if h == 12 {
		tm = "pm"
	}

	hs := fmt.Sprint(h)
	ms := fmt.Sprint(m)

	if m < 10 && m != 0 {
		ms = "0" + ms
	}
	return hs + ":" + ms + " " + tm
}

func LoadSchedules() []Schedule {
	file, err := os.Open("../../Schedules/Schedules.csv")

	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		log.Println(err)
	}

	var s []Schedule

	for _, r := range records[1:] {
		f := strings.Split(r[3],"|")
		var fs []int
		for _,x := range f {
			cf,err := strconv.Atoi(strings.Trim(x," "))
			if err != nil {
				log.Println(err)
			}
			fs = append(fs, cf)
		}
		ns := NewSchedule(r[1], r[2], r[0],fs...)
		s = append(s, ns)
	}
	return s
}
