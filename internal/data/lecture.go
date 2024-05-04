package data

import (
	"gorm.io/gorm"
	"talk_rater_bot/internal/validator"
	"time"
)

type Lecture struct {
	ID      int64 // ID parsed from URL
	Name    string
	Speaker string
	URL     string    `gorm:"unique;index"`
	Start   time.Time `gorm:"index"`
	End     time.Time

	UsersInFavourite []*User        `gorm:"many2many:user_lectures;"` // Only for gorm. Many2Many
	DeletedAt        gorm.DeletedAt `gorm:"index"`
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
