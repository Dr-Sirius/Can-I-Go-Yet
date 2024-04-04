package scheduler

import "time"


const (
	OPEN int = iota
	
)


type Schedule struct {
	StartTime time.Time
	EndTime time.Time
	flags   map[int]bool
}