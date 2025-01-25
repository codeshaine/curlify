package request

import (
	"testing"
)

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      string
	}{
		{"", "", "empty URL"},
		{"example.com", "", "missing scheme"},
		{"ftp://example.com", "", "only http and https allowed"},
		{"http://", "", "missing host"},
		{"http://example.com", "http://example.com", ""},
		{"https://example.com", "https://example.com", ""},
	}

	for _, test := range tests {
		result, err := ValidateUrl(test.input)
		if result != test.expected {
			t.Errorf("ValidateUrl(%q) = %q; want %q", test.input, result, test.expected)
		}
		if err != nil && err.Error() != test.err {
			t.Errorf("ValidateUrl(%q) error = %q; want %q", test.input, err.Error(), test.err)
		}
		if err == nil && test.err != "" {
			t.Errorf("ValidateUrl(%q) expected error %q; got nil", test.input, test.err)
		}
	}
}
