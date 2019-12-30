package caltrain

// timeTable keeps track of the day's timetable. Should only need to be called
// once per day
type timeTable struct{}

func newTimeTable() *timeTable {
	return &timeTable{}
}
