package validators

import "testing"

func TestValidator(t *testing.T) {
	tests := []struct {
		name   string
		checks []struct {
			ok      bool
			key     string
			message string
		}
		expectValid    bool
		expectedErrors map[string]string
	}{
		{
			name: "All checks pass",
			checks: []struct {
				ok      bool
				key     string
				message string
			}{
				{ok: true, key: "title", message: "Title is required"},
				{ok: true, key: "speaker", message: "Speaker is required"},
			},
			expectValid:    true,
			expectedErrors: map[string]string{},
		},
		{
			name: "One check fails",
			checks: []struct {
				ok      bool
				key     string
				message string
			}{
				{ok: true, key: "title", message: "Title is required"},
				{ok: false, key: "speaker", message: "Speaker is required"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"speaker": "Speaker is required"},
		},
		{
			name: "Multiple checks fail",
			checks: []struct {
				ok      bool
				key     string
				message string
			}{
				{ok: false, key: "title", message: "Title is required"},
				{ok: false, key: "speaker", message: "Speaker is required"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"title": "Title is required", "speaker": "Speaker is required"},
		},
		{
			name: "Duplicate error keys",
			checks: []struct {
				ok      bool
				key     string
				message string
			}{
				{ok: false, key: "title", message: "Title is required"},
				{ok: false, key: "title", message: "Title is still required"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"title": "Title is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			for _, check := range tt.checks {
				v.Check(check.ok, check.key, check.message)
			}
			if v.Valid() != tt.expectValid {
				t.Errorf("Expected valid: %v, got: %v", tt.expectValid, v.Valid())
			}
			if !equalErrorMaps(v.Errors, tt.expectedErrors) {
				t.Errorf("Expected errors: %v, got: %v", tt.expectedErrors, v.Errors)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected bool
	}{
		{
			name:     "All unique",
			values:   []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "Duplicates present",
			values:   []string{"a", "b", "a"},
			expected: false,
		},
		{
			name:     "Empty slice",
			values:   []string{},
			expected: true,
		},
		{
			name:     "Single element",
			values:   []string{"a"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.values)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
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
