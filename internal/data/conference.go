package data

import (
	"TalkRater_Bot/internal/validator"
	"time"
)

type Conference struct {
	Name              string    `json:"name"`
	URL               string    `json:"url"`
	StartTime         time.Time `json:"start"`
	EndTime           time.Time `json:"end"`
	EndEvaluationTime time.Time `json:"end_evaluation"`
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
