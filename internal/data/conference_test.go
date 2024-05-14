package data

import (
	"talk_rater_bot/internal/validators"
	"testing"
	"time"
)

func TestAreEqualConferences(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(2 * time.Hour)
	t3 := t2.Add(1 * time.Hour)

	tests := []struct {
		name     string
		c1       *Conference
		c2       *Conference
		expected bool
	}{
		{
			name:     "both nil",
			c1:       nil,
			c2:       nil,
			expected: true,
		},
		{
			name:     "one nil",
			c1:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			c2:       nil,
			expected: false,
		},
		{
			name:     "different names",
			c1:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			c2:       &Conference{Name: "Conf B", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			expected: false,
		},
		{
			name:     "different URLs",
			c1:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			c2:       &Conference{Name: "Conf A", URL: "http://confB.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			expected: false,
		},
		{
			name:     "different start times",
			c1:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			c2:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1.Add(1 * time.Minute), EndTime: t2, EndEvaluationTime: t3},
			expected: false,
		},
		{
			name:     "different end times",
			c1:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			c2:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2.Add(1 * time.Minute), EndEvaluationTime: t3},
			expected: false,
		},
		{
			name:     "different end evaluation times",
			c1:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			c2:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3.Add(1 * time.Minute)},
			expected: false,
		},
		{
			name:     "all fields equal",
			c1:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			c2:       &Conference{Name: "Conf A", URL: "http://confA.com", StartTime: t1, EndTime: t2, EndEvaluationTime: t3},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AreEqualConferences(tt.c1, tt.c2); got != tt.expected {
				t.Errorf("AreEqualConferences() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateConference(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(2 * time.Hour)
	t3 := t2.Add(1 * time.Hour)

	tests := []struct {
		name           string
		conference     Conference
		expectedErrors map[string]string
	}{
		{
			name: "valid conference",
			conference: Conference{
				Name:              "Conf A",
				URL:               "http://confA.com",
				StartTime:         t1,
				EndTime:           t2,
				EndEvaluationTime: t3,
			},
			expectedErrors: map[string]string{},
		},
		{
			name: "missing name",
			conference: Conference{
				Name:              "",
				URL:               "http://confA.com",
				StartTime:         t1,
				EndTime:           t2,
				EndEvaluationTime: t3,
			},
			expectedErrors: map[string]string{"name": "Title is required"},
		},
		{
			name: "missing URL",
			conference: Conference{
				Name:              "Conf A",
				URL:               "",
				StartTime:         t1,
				EndTime:           t2,
				EndEvaluationTime: t3,
			},
			expectedErrors: map[string]string{"url": "URL is required"},
		},
		{
			name: "missing start time",
			conference: Conference{
				Name:              "Conf A",
				URL:               "http://confA.com",
				StartTime:         time.Time{},
				EndTime:           t2,
				EndEvaluationTime: t3,
			},
			expectedErrors: map[string]string{"start": "Start time is required"},
		},
		{
			name: "missing end time",
			conference: Conference{
				Name:              "Conf A",
				URL:               "http://confA.com",
				StartTime:         t1,
				EndTime:           time.Time{},
				EndEvaluationTime: t3,
			},
			expectedErrors: map[string]string{"start and end time": "Start time must be before End time"},
		},
		{
			name: "missing end evaluation time",
			conference: Conference{
				Name:              "Conf A",
				URL:               "http://confA.com",
				StartTime:         t1,
				EndTime:           t2,
				EndEvaluationTime: time.Time{},
			},
			expectedErrors: map[string]string{"end and end evaluation time": "End time must be before End evaluation time"},
		},
		{
			name: "start time after end time",
			conference: Conference{
				Name:              "Conf A",
				URL:               "http://confA.com",
				StartTime:         t2,
				EndTime:           t1,
				EndEvaluationTime: t3,
			},
			expectedErrors: map[string]string{"start and end time": "Start time must be before End time"},
		},
		{
			name: "end time after end evaluation time",
			conference: Conference{
				Name:              "Conf A",
				URL:               "http://confA.com",
				StartTime:         t1,
				EndTime:           t3,
				EndEvaluationTime: t2,
			},
			expectedErrors: map[string]string{"end and end evaluation time": "End time must be before End evaluation time"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validators.New()
			ValidateConference(v, &tt.conference)

			if len(v.Errors) != len(tt.expectedErrors) {
				t.Errorf("expected %d v.Errors, got %d. Errors: %+v", len(tt.expectedErrors), len(v.Errors), v.Errors)
			}

			for key, expectedMessage := range tt.expectedErrors {
				if message, exists := v.Errors[key]; !exists || message != expectedMessage {
					t.Errorf("expected error for %s: %s, got %s", key, expectedMessage, message)
				}
			}
		})
	}
}
