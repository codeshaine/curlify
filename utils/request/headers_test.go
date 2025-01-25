package request

import (
	"testing"
)

func TestParseHeaders(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedErr   bool
		expectedCount int
	}{
		{
			name:          "Valid Single Header",
			input:         "Content-Type: application/json",
			expectedErr:   false,
			expectedCount: 1,
		},
		{
			name:          "Multiple Valid Headers",
			input:         "Content-Type: application/json\nAuthorization: Bearer token123",
			expectedErr:   false,
			expectedCount: 2,
		},
		{
			name:          "Empty Input",
			input:         "",
			expectedErr:   false,
			expectedCount: 0,
		},
		{
			name:          "Missing Colon",
			input:         "Invalid-Header",
			expectedErr:   true,
			expectedCount: 0,
		},
		{
			name:          "Empty Header Name",
			input:         ": value",
			expectedErr:   true,
			expectedCount: 0,
		},
		{
			name:          "Invalid Header Name",
			input:         "Invalid Header!: value",
			expectedErr:   true,
			expectedCount: 0,
		},
		{
			name:          "Empty Header Value",
			input:         "Content-Type:",
			expectedErr:   true,
			expectedCount: 0,
		},
		{
			name:          "Multiple Values for Same Header",
			input:         "X-Custom: value1\nX-Custom: value2",
			expectedErr:   false,
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers, err := ParseHeaders(tc.input)

			if tc.expectedErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(headers) != tc.expectedCount {
					t.Errorf("Expected %d headers, got %d", tc.expectedCount, len(headers))
				}
			}
		})
	}
}

func TestValidateHeaderName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid Header Name", "Content-Type", true},
		{"Header Name with Numbers", "X-Custom-Header-123", true},
		{"Invalid Header Name", "Invalid Header!", false},
		{"Empty Header Name", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateHeaderName(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v for input %s, got %v", tc.expected, tc.input, result)
			}
		})
	}
}
