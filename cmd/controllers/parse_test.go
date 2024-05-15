package controllers

import (
	"os"
	"testing"

	"gorm.io/gorm"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/helpers"
)

func createTempCSV(content string) (string, error) {
	tmpfile, err := os.CreateTemp("", "lectures-*.csv")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		return "", err
	}

	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func TestParseSchedule(t *testing.T) {
	csvContent := `StartTime,Duration,Title,Speaker,URL
21/07/2024 10:00:00,60,Go Concurrency Patterns,John Doe,https://example.com/go-concurrency
21/07/2024 11:00:00,45,Building Microservices,Jane Smith,https://example.com/microservices
`

	csvFilePath, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(csvFilePath)

	timeParser, _ := helpers.NewTimeParserInMoscow()
	controller := &Controller{
		db:         &gorm.DB{},
		timeParser: timeParser,
		conference: &data.Conference{},
	}

	// Call parseSchedule method
	lectures, err := controller.parseSchedule(csvFilePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the number of parsed lectures
	if len(lectures) != 2 {
		t.Fatalf("Expected 2 lectures, got %d", len(lectures))
	}

	// Verify the first lecture
	if lectures[0].StartTime != "21/07/2024 10:00:00" {
		t.Errorf("Expected StartTime '21/07/2024 10:00:00', got %s", lectures[0].StartTime)
	}
	if lectures[0].Duration != "60m" {
		t.Errorf("Expected Duration '60m', got %s", lectures[0].Duration)
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

	// Verify the second lecture
	if lectures[1].StartTime != "21/07/2024 11:00:00" {
		t.Errorf("Expected StartTime '21/07/2024 11:00:00', got %s", lectures[1].StartTime)
	}
	if lectures[1].Duration != "45m" {
		t.Errorf("Expected Duration '45m', got %s", lectures[1].Duration)
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
}

func TestParseSchedule_InvalidCSV(t *testing.T) {
	invalidCSVContent := `StartTime,Duration,Title,Speaker
21/07/2024 10:00:00,60,Go Concurrency Patterns,John Doe
`

	csvFilePath, err := createTempCSV(invalidCSVContent)
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(csvFilePath)

	timeParser, _ := helpers.NewTimeParserInMoscow()

	controller := &Controller{
		db:         &gorm.DB{},
		timeParser: timeParser,
		conference: &data.Conference{},
	}

	_, err = controller.parseSchedule(csvFilePath)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	expectedErr := "count of args should be 5, got 4.\nrecord: [21/07/2024 10:00:00 60 Go Concurrency Patterns John Doe]"
	if err.Error() != expectedErr {
		t.Fatalf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
