package scheduler

import (
	"log"
	"testing"
)

func TestStringStartTime(t *testing.T) {
	s := NewSchedule("7:01 am", "1:02 pm", "2024-08-19", 0)
	log.Println(s.StringStartTime(), s.StringEndTime())
}



func TestNewDay(t *testing.T) {
	s := []Schedule{
		NewSchedule("2:03 pm", "2:22 pm", "2024-08-19", 0),
		NewSchedule("7:01 am", "1:02 pm", "2024-07-19", 0),
		NewSchedule("1:01 pm", "2:02 pm", "2024-08-19", 0),
		NewSchedule("7:01 am", "1:02 pm", "2024-07-19", 0),
		NewSchedule("7:09 am", "12:05 pm", "2024-08-19", 0),
		
	}
	log.Println(NewDay("2024-08-19",s...))
}