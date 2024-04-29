package data

import "time"

type Conference struct {
	Name              string    `json:"name"`
	URL               string    `json:"url"`
	StartTime         time.Time `json:"start"`
	EndTime           time.Time `json:"end"`
	EndEvaluationTime time.Time `json:"end_evaluation"`
}
