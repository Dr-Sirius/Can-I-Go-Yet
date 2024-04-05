package scheduler

import (
	"log"
	"testing"
)

func TestStringStartTime(t *testing.T) {
	s := NewSchedule("7:01 am", "1:02 pm", "2024-08-19", 0)
	log.Println(s.StringStartTime(), s.StringEndTime())
}
