package controllers

import (
	"talk_rater_bot/internal/helpers"
	"testing"
	"time"

	"talk_rater_bot/internal/data"
)

func TestConvertAndValidateLectures(t *testing.T) {
	// Set up mock data
	lecturesInput := []*LectureInput{
		{
			StartTime: "21/07/2024 10:00:00",
			Duration:  "60m",
			Title:     "Go Concurrency Patterns",
			Speaker:   "John Doe",
			URL:       "https://example.com/go-concurrency",
		},
		{
			StartTime: "21/07/2024 11:00:00",
			Duration:  "45m",
			Title:     "Building Microservices",
			Speaker:   "Jane Smith",
			URL:       "https://example.com/microservices",
		},
	}

	conference := &data.Conference{
		StartTime: time.Date(2024, 6, 21, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 8, 21, 18, 0, 0, 0, time.UTC),
	}

	timeParser, _ := helpers.NewTimeParserInMoscow()
	controller := &Controller{
		db:         nil,
		timeParser: timeParser,
		conference: conference,
	}

	lectures, err := controller.convertAndValidateLectures(lecturesInput)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(lectures) != 2 {
		t.Fatalf("Expected 2 lectures, got %d", len(lectures))
	}

	if lectures[0].Title != "Go Concurrency Patterns" {
		t.Errorf("Expected Title 'Go Concurrency Patterns', got %s", lectures[0].Title)
	}
	if lectures[0].Speaker != "John Doe" {
		t.Errorf("Expected Speaker 'John Doe', got %s", lectures[0].Speaker)
	}
	if lectures[0].URL != "https://example.com/go-concurrency" {
		t.Errorf("Expected URL 'https://example.com/go-concurrency', got %s", lectures[0].URL)
	}
	if lectures[0].Start.Equal(time.Date(2024, 7, 21, 8, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected Start '21/07/2024 10:00:00', got %s", lectures[0].Start)
	}
	if lectures[0].End.Equal(time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected End '21/07/2024 11:00:00', got %s", lectures[0].End)
	}

	if lectures[1].Title != "Building Microservices" {
		t.Errorf("Expected Title 'Building Microservices', got %s", lectures[1].Title)
	}
	if lectures[1].Speaker != "Jane Smith" {
		t.Errorf("Expected Speaker 'Jane Smith', got %s", lectures[1].Speaker)
	}
	if lectures[1].URL != "https://example.com/microservices" {
		t.Errorf("Expected URL 'https://example.com/microservices', got %s", lectures[1].URL)
	}
	if lectures[1].Start.Equal(time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected Start '21/07/2024 11:00:00', got %s", lectures[1].Start)
	}
	if lectures[1].End.Equal(time.Date(2024, 7, 21, 9, 45, 0, 0, time.UTC)) {
		t.Errorf("Expected End '21/07/2024 11:45:00', got %s", lectures[1].End)
	}
}

func TestConvertAndValidateLectures_DuplicateURL(t *testing.T) {
	lecturesInput := []*LectureInput{
		{
			StartTime: "21/07/2024 10:00:00",
			Duration:  "60m",
			Title:     "Go Concurrency Patterns",
			Speaker:   "John Doe",
			URL:       "https://example.com/go-concurrency",
		},
		{
			StartTime: "21/07/2024 11:00:00",
			Duration:  "45m",
			Title:     "Building Microservices",
			Speaker:   "Jane Smith",
			URL:       "https://example.com/go-concurrency", // Duplicate URL
		},
	}

	// Set up conference data
	conference := &data.Conference{
		StartTime: time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 7, 21, 18, 0, 0, 0, time.UTC),
	}

	timeParser, _ := helpers.NewTimeParserInMoscow()
	controller := &Controller{
		db:         nil,
		timeParser: timeParser,
		conference: conference,
	}

	_, err := controller.convertAndValidateLectures(lecturesInput)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	expectedErr := "URLs must be unique"
	if err.Error() != expectedErr {
		t.Fatalf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestConvertAndValidateLectures_InvalidTime(t *testing.T) {
	lecturesInput := []*LectureInput{
		{
			StartTime: "invalid time",
			Duration:  "60m",
			Title:     "Go Concurrency Patterns",
			Speaker:   "John Doe",
			URL:       "https://example.com/go-concurrency",
		},
		{
			StartTime: "21/07/2024 11:00:00",
			Duration:  "45m",
			Title:     "Building Microservices",
			Speaker:   "Jane Smith",
			URL:       "https://example.com/microservices",
		},
	}

	conference := &data.Conference{
		StartTime: time.Date(2024, 7, 20, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 7, 22, 18, 0, 0, 0, time.UTC),
	}

	timeParser, _ := helpers.NewTimeParserInMoscow()
	controller := &Controller{
		db:         nil,
		timeParser: timeParser,
		conference: conference,
	}

	_, err := controller.convertAndValidateLectures(lecturesInput)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	expectedErr := "error in line: &{StartTime:invalid time Duration:60m Title:Go Concurrency Patterns Speaker:John Doe URL:https://example.com/go-concurrency}, description: parsing time \"invalid time\" as \"02/01/2006 15:04:05\": cannot parse \"invalid time\" as \"02\""
	if err.Error() != expectedErr {
		t.Fatalf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
