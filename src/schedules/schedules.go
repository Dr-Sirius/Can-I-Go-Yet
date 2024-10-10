package schedules

import (
	

	"fmt"
	"log"

	"slices"
	
	"time"
)

// Flags
const (
	OPEN int = iota
	BRKE     // break
	UNDS     // understaffed
	HDAY     // holiday
)

// A struct for holding schedule info
type Schedule struct {
	StartTime time.Time
	EndTime   time.Time
	Flags     []int
}
// variables
var timeFormat string = "2006-01-02 3:04 pm"

/* 
-------------------------------Schedule Creation---------------------------------------
*/

/*
Creates a new Schedule struct from formated strings and flags

EX. st & et -> 12:30 am - date -> 2024-08-19
*/
func New(st string, et string, date string, flags []int) Schedule {
	slices.Sort(flags)
	return NewFromTime(convertTime(st, date), convertTime(et, date), flags)
}

/*
Creates a new Schedule struct from time.Time structs and flags
*/
func NewFromTime(st time.Time, et time.Time, flags []int) Schedule {
	return Schedule{st, et, flags}
}


/*
Sorts schedules based on the start time

Uses a bubble sort algorithm
*/
func scheduleSort(Schedules []Schedule) []Schedule {
	for i := len(Schedules) - 1; i >= 0; i -= 1 {

		for x := range i {
			if Schedules[x].EqualTimes(Schedules[x+1]) {
				newSchedules := append(Schedules[:x],Schedules[x+1:]...)
				return scheduleSort(newSchedules)
			} else if Schedules[x].StartTime.Compare(Schedules[x+1].StartTime) > 0 {
				temp := Schedules[x]
				Schedules[x] = Schedules[x+1]
				Schedules[x+1] = temp
			}

		}
	}
	return Schedules
}

/* 
-------------------------------String Methods---------------------------------------
*/

/*
Returns string version of schedule
*/
func (s Schedule) String() string {
	st := s.StringStartTime()
	et := s.StringEndTime()
	return st + " - " + et + " " + fmt.Sprint(s.Flags)
}

/*
Returns formated string version of schedule
*/
func (s Schedule) PrettyString() string {
	st := s.StringStartTime()
	et := s.StringEndTime()
	f := func() string {
		str := ""
		if s.Flags[0] == -1 {
			return "|Closed"
		}
		for i := range s.Flags {
			if i == OPEN {
				str += "|Open"
			}
			if i == BRKE {
				str += "|Break"
			}
			if i == UNDS {
				str += "|Understaffed"
			}
			if i == HDAY {
				str += "|Holiday"
			}
			
		}
		return str
	}
	return st + " - " + et + " " + f()
}

/*
Returns the StartTime as a formated string

EX. time.Time -> 7:00 am
*/
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

	if m < 10 {
		ms = "0" + ms
	}
	return hs + ":" + ms + " " + tm
}

/*
Returns the EndTime as a formated string

EX. time.Time -> 7:00 am
*/
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

	if m < 10 {
		ms = "0" + ms
	}
	return hs + ":" + ms + " " + tm
}

/*
Returns string containg date formatted as yyyy-mm-dd
*/
func (s Schedule) Date() string {
	y, m, d := s.StartTime.Date()
	mt := fmt.Sprint(int(m))
	if int(m) < 10 {
		mt = "0" + mt
	}

	return fmt.Sprint(y) + "-" + mt + "-" + fmt.Sprint(d)
}

/*
converts string formatted dates into a time.Time struct
*/
func convertTime(timeString string, dateString string) time.Time {
	timeString = dateString + " " + timeString
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println(err)
	}
	timeStruct, err := time.ParseInLocation(timeFormat, timeString, loc)
	if err != nil {
		log.Println(err)
	}
	return timeStruct
}

/* 
-------------------------------Comparison Methods---------------------------------------
*/

/*
Returns true if schedule s is equal to schedule o
*/
func (s Schedule) Equal(o Schedule) bool {
	return s.StartTime.Equal(o.StartTime) && s.EndTime.Equal(o.EndTime) && func() bool {
		
		if len(s.Flags) != len(o.Flags) {
			return false
		}
		
		for a := range len(s.Flags) {
			if s.Flags[a] != o.Flags[a] {
				return false
			}
		}
		return true

	}()
}

/*
Returns true if schedule s StartTime and EndTime is equal to schedule o StartTime and EndTime
*/
func (s Schedule) EqualTimes(o Schedule) bool {
	return s.StartTime.Equal(o.StartTime) && s.EndTime.Equal(o.EndTime)
}


// /* 
// -------------------------------File Methods---------------------------------------
// */

// /*
// Returns an array of Schedule structs from the Schedules folder
// */
// func LoadSchedules() []Schedule {
// 	if _, err := os.Stat("Schedules/Schedules.csv"); errors.Is(err, os.ErrNotExist) {
// 		if err := os.Mkdir("Schedules", os.ModePerm); err != nil {
// 			log.Println(err)
// 		} else {
// 			os.Create("Schedules/Schedules.csv")
// 			os.WriteFile("Schedules/Schedules.csv", []byte("Date, Start_Time, End_Time, Flags"), os.ModePerm)
// 		}
// 	}
// 	file, err := os.Open("Schedules/Schedules.csv")

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	defer file.Close()

// 	reader := csv.NewReader(file)

// 	records, err := reader.ReadAll()

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var s []Schedule

// 	for _, r := range records[1:] {
// 		f := strings.Split(r[3], "|")
// 		var fs []int
// 		for _, x := range f {
// 			cf, err := strconv.Atoi(strings.Trim(x, " "))
// 			if !errors.Is(err,strconv.ErrSyntax) && err != nil {
// 				log.Println(err)
// 			}
// 			fs = append(fs, cf)
// 		}
// 		ns := New(r[1], r[2], r[0], fs)
// 		s = append(s, ns)
// 	}
// 	return s
// }

// /*
// Adds schedule to Schedules.csv
// */
// func AddSchedule(date string, startTime string, endTime string, flags []int) {
// 	flagString := ""
// 	for _, f := range flags {
// 		flagString += fmt.Sprint(f) + "|"
// 	}

// 	sch := "\n" + date + ", " + startTime + ", " + endTime + ", " + flagString

// 	file, err := os.OpenFile("Schedules/Schedules.csv", os.O_APPEND|os.O_WRONLY, 0644)

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	defer file.Close()

// 	_, err = file.WriteString(sch)

// 	if err != nil {
// 		log.Println(err)
		

// 	}
// }


// /*
// Removes passed schedule from Schedules.csv
// */
// func RemoveSchedule(s Schedule) {
// 	sch := LoadSchedules()
// 	for i, x := range sch {
// 		if x.Equal(s) {
// 			sch = append(sch[0:i], sch[i+1:]...)
// 			break
// 		}
// 	}
// 	if err := os.WriteFile("Schedules/Schedules.csv",[]byte("Date, Start_Time, End_Time, Flags"),os.ModePerm); err != nil {
// 		log.Println(err)
// 	}
	

// 	for _,x := range sch {
// 		AddSchedule(x.Date(),x.StringStartTime(),x.StringEndTime(),x.Flags)
// 	}
// }
