package helpers

import (
	"testing"
	"time"
)

func TestTimeParser_ParseTimeString(t *testing.T) {
	timeString := "14/05/2024 15:00:00" // in MSK time zone
	expectedTime := time.Date(2024, 5, 14, 12, 0, 0, 0, time.UTC)

	tp, err := NewTimeParserInMoscow()
	if err != nil {
		t.Fatalf("Failed to create TimeParser: %s", err)
	}

	parsedTime, err := tp.ParseTimeString(timeString)
	if err != nil {
		t.Fatalf("Failed to parse time string: %s", err)
	}

	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Parsed time does not match expected time. Expected: %s, Got: %s", expectedTime, parsedTime)
	}
}

func TestTimeParser_ConvertTime(t *testing.T) {
	expectedTimeString := "14/05/2024 15:00:00"
	inputTime := time.Date(2024, 5, 14, 12, 0, 0, 0, time.UTC)

	tp, err := NewTimeParserInMoscow()
	if err != nil {
		t.Fatalf("Failed to create TimeParser: %s", err)
	}

	convertedTime := tp.ConvertTime(inputTime)
	if convertedTime != expectedTimeString {
		t.Errorf("Converted time string does not match expected. Expected: %s, Got: %s", expectedTimeString, convertedTime)
	}
}
