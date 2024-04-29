package data

import (
	"TalkRater_Bot/internal/validator"
	"time"
)

type Lecture struct {
	Id      int64
	Name    string
	Speaker string
	URL     string //check contains unique
	Start   time.Time
	End     time.Time
}

func (l *Lecture) Duration() time.Duration {
	return l.End.Sub(l.Start)
}

func ValidateLecture(v *validator.Validator, lect *Lecture) {
	v.Check(lect.Name != "", "name", "Name is required")
	v.Check(lect.Speaker != "", "speaker", "Speaker is required")
	v.Check(!lect.Start.IsZero(), "start", "Start time is required")
	v.Check(!lect.End.IsZero(), "end", "End time is required")

	v.Check(lect.Start.Before(lect.End), "start and end time", "Start time must be before End time")

	diff := lect.End.Sub(lect.Start)
	v.Check(diff <= time.Hour*12, "duration", "Duration must be lower than 12 hours")
}
