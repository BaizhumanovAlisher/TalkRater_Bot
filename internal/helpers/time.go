package helpers

import (
	"time"
)

type TimeParser struct {
	FileLayout string
	Location   *time.Location
}

func NewTimeParserInMoscow() (*TimeParser, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	return &TimeParser{
		FileLayout: "02/01/2006 15:04:05",
		Location:   location,
	}, nil
}

func (tp *TimeParser) ParseTimeString(timeString string) (time.Time, error) {
	t, err := time.ParseInLocation(tp.FileLayout, timeString, tp.Location)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func (tp *TimeParser) ConvertTime(t time.Time) string {
	return t.In(tp.Location).Format(tp.FileLayout)
}
