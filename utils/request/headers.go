package request

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	validHeaderNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
)

func ParseHeaders(input string) (map[string][]string, error) {
	headers := make(map[string][]string)

	if input == "" {
		return headers, nil
	}

	lines := strings.Split(input, "\n")
	for lineNum, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format at line %d: missing colon", lineNum+1)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty header name at line %d", lineNum+1)
		}

		if !validHeaderNameRegex.MatchString(key) {
			return nil, fmt.Errorf("invalid header name at line %d: %s", lineNum+1, key)
		}

		if value == "" {
			return nil, fmt.Errorf("empty header value for %s at line %d", key, lineNum+1)
		}

		headers[key] = append(headers[key], value)
	}

	return headers, nil
}

func ValidateHeaderName(name string) bool {
	return validHeaderNameRegex.MatchString(name)
}
