package scheduler

import (
	"fmt"
	"log"
	"time"
)

const (
	OPEN int = iota
	UNDS     // understaffed
	HDAY     // holiday
)

type Schedule struct {
	Date 	  string
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
	return NewScheduleFromTime(date, convertTime(st, date), convertTime(et, date), flags...)
}

func NewScheduleFromTime(date string, st time.Time, et time.Time, flags ...int) Schedule {
	f := make(map[int]bool)
	for _, x := range flags {
		f[x] = true
	}
	return Schedule{date,st, et, f}
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
	s := make([]Schedule,0)
	for _,x := range schedules {
		if x.Date == date {
			s = append(s, x)
		}
	}
	return Day{date,s}
}


func (s Schedule) StringStartTime() string {
	h,m,_ := s.StartTime.Clock()
	tm := "am"

	if h > 12 {
		h = h - 12
		tm = "pm"
	}
	if h == 12 {tm = "pm"}

	hs := fmt.Sprint(h)
	ms := fmt.Sprint(m)
	
	if m < 10 && m!=0 {
		ms = "0" + ms
	}
	return hs + ":" + ms + " " + tm 
}

func (s Schedule) StringEndTime() string {
	h,m,_ := s.EndTime.Clock()
	tm := "am"

	if h > 12 {
		h = h - 12
		tm = "pm"
	}
	if h == 12 {tm = "pm"}

	hs := fmt.Sprint(h)
	ms := fmt.Sprint(m)
	
	if m < 10 && m!=0 {
		ms = "0" + ms
	}
	return hs + ":" + ms + " " + tm 
}