package data

import (
	"talk_rater_bot/internal/validators"
	"testing"
)

func TestValidateEvaluation(t *testing.T) {
	tests := []struct {
		name           string
		eval           Evaluation
		expectedErrors map[string]string
	}{
		{
			name: "valid evaluation",
			eval: Evaluation{
				UserID:           1,
				LectureID:        1,
				TypeEvaluation:   Correct,
				ScoreContent:     3,
				ScorePerformance: 4,
			},
			expectedErrors: map[string]string{},
		},
		{
			name: "missing user ID",
			eval: Evaluation{
				UserID:           0,
				LectureID:        1,
				TypeEvaluation:   Correct,
				ScoreContent:     3,
				ScorePerformance: 4,
			},
			expectedErrors: map[string]string{"user's id": "user's id должен быть больше нуля"},
		},
		{
			name: "missing lecture ID",
			eval: Evaluation{
				UserID:           1,
				LectureID:        0,
				TypeEvaluation:   Correct,
				ScoreContent:     3,
				ScorePerformance: 4,
			},
			expectedErrors: map[string]string{"lecture's id": "lecture's id должен быть больше нуля"},
		},
		{
			name: "invalid score content",
			eval: Evaluation{
				UserID:           1,
				LectureID:        1,
				TypeEvaluation:   Correct,
				ScoreContent:     6,
				ScorePerformance: 4,
			},
			expectedErrors: map[string]string{"score content": "оценка содержания должна быть между 1 и 5"},
		},
		{
			name: "invalid score performance",
			eval: Evaluation{
				UserID:           1,
				LectureID:        1,
				TypeEvaluation:   Correct,
				ScoreContent:     3,
				ScorePerformance: 0,
			},
			expectedErrors: map[string]string{"score performance": "оценка выступления должна быть между 1 и 5"},
		},
		{
			name: "non-correct type evaluation, no score check",
			eval: Evaluation{
				UserID:           1,
				LectureID:        1,
				TypeEvaluation:   NoEvaluation,
				ScoreContent:     0,
				ScorePerformance: 0,
			},
			expectedErrors: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validators.New()
			ValidateEvaluation(v, &tt.eval)
			errors := v.Errors

			if len(errors) != len(tt.expectedErrors) {
				t.Errorf("expected %d errors, got %d", len(tt.expectedErrors), len(errors))
			}

			for key, expectedMessage := range tt.expectedErrors {
				if message, exists := errors[key]; !exists || message != expectedMessage {
					t.Errorf("expected error for %s: %s, got %s", key, expectedMessage, message)
				}
			}
		})
	}
}
