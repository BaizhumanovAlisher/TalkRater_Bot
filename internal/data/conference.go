package data

import (
	"talk_rater_bot/internal/validator"
	"time"
)

type Conference struct {
	Name              string `gorm:"primaryKey"`
	URL               string
	StartTime         time.Time
	EndTime           time.Time
	EndEvaluationTime time.Time
}

func ValidateConference(v *validator.Validator, conf *Conference) {
	v.Check(conf.Name != "", "name", "Name is required")
	v.Check(conf.URL != "", "url", "URL is required")

	v.Check(!conf.StartTime.IsZero(), "start", "Start time is required")
	v.Check(!conf.EndTime.IsZero(), "end", "End time is required")
	v.Check(!conf.EndEvaluationTime.IsZero(), "end evaluation", "End evaluation time is required")

	v.Check(conf.StartTime.Before(conf.EndTime), "start and end time", "Start time must be before End time")
	v.Check(conf.EndTime.Before(conf.EndEvaluationTime), "end and end evaluation time", "End time must be before End evaluation time")
}

const (
	FileLayout = "02/01/2006 15:04:05"
)

func ParseTimeString(timeString string, location *time.Location, layout string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, timeString, location)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func ConvertTime(t time.Time, location *time.Location, layout string) string {
	return t.In(location).Format(layout)
}

func AreEqualConferences(c1 *Conference, c2 *Conference) bool {
	if c1 == c2 {
		return true
	}

	if c1 == nil || c2 == nil {
		return false
	}

	return c1.Name == c2.Name && c1.URL == c2.URL &&
		c1.StartTime.Equal(c2.StartTime) &&
		c1.EndTime.Equal(c2.EndTime) &&
		c1.EndEvaluationTime.Equal(c2.EndEvaluationTime)
}
