package request

import (
	"fmt"
	"strings"
)

func ValidateUrl(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("empty URL")
	}

	// Check for scheme
	schemeParts := strings.SplitN(input, "://", 2)
	if len(schemeParts) != 2 {
		return "", fmt.Errorf("invalid url")
	}

	// Validate scheme
	scheme := schemeParts[0]
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("only http and https allowed")
	}

	// Validate host
	hostPath := strings.SplitN(schemeParts[1], "/", 2)
	if hostPath[0] == "" {
		return "", fmt.Errorf("missing host")
	}

	return input, nil
}
