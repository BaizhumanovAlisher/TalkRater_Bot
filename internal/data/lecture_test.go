package data

import (
	"talk_rater_bot/internal/validators"
	"testing"
	"time"
)

func TestValidateLecture(t *testing.T) {
	tests := []struct {
		name        string
		lecture     *Lecture
		expectValid bool
		expectedErr map[string]string
	}{
		{
			name: "All fields are valid",
			lecture: &Lecture{
				Title:   "Go Concurrency Patterns",
				Speaker: "John Doe",
				Start:   time.Now().Add(time.Hour),
				End:     time.Now().Add(2 * time.Hour),
			},
			expectValid: true,
		},
		{
			name: "Missing title",
			lecture: &Lecture{
				Title:   "",
				Speaker: "John Doe",
				Start:   time.Now().Add(time.Hour),
				End:     time.Now().Add(2 * time.Hour),
			},
			expectValid: false,
			expectedErr: map[string]string{"name": "Title is required"},
		},
		{
			name: "Missing speaker",
			lecture: &Lecture{
				Title:   "Go Concurrency Patterns",
				Speaker: "",
				Start:   time.Now().Add(time.Hour),
				End:     time.Now().Add(2 * time.Hour),
			},
			expectValid: false,
			expectedErr: map[string]string{"speaker": "Speaker is required"},
		},
		{
			name: "Missing start time",
			lecture: &Lecture{
				Title:   "Go Concurrency Patterns",
				Speaker: "John Doe",
				Start:   time.Time{},
				End:     time.Time{}.Add(time.Hour),
			},
			expectValid: false,
			expectedErr: map[string]string{"start": "Start time is required"},
		},
		{
			name: "Start time is after end time",
			lecture: &Lecture{
				Title:   "Go Concurrency Patterns",
				Speaker: "John Doe",
				Start:   time.Now().Add(2 * time.Hour),
				End:     time.Now().Add(time.Hour),
			},
			expectValid: false,
			expectedErr: map[string]string{"start and end time": "Start time must be before End time"},
		},
		{
			name: "Duration longer than 12 hours",
			lecture: &Lecture{
				Title:   "Go Concurrency Patterns",
				Speaker: "John Doe",
				Start:   time.Now(),
				End:     time.Now().Add(13 * time.Hour),
			},
			expectValid: false,
			expectedErr: map[string]string{"duration": "Duration must be lower than 12 hours"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := validators.New()
			ValidateLecture(validator, tt.lecture)

			if validator.Valid() != tt.expectValid {
				t.Errorf("Expected valid: %v, got: %v", tt.expectValid, validator.Valid())
			}

			if !validator.Valid() && !equalErrorMaps(validator.Errors, tt.expectedErr) {
				t.Errorf("Expected errors: %v, got: %v", tt.expectedErr, validator.Errors)
			}
		})
	}
}

func equalErrorMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
