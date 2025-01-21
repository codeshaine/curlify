package request

import "strings"

func ParseHeaders(input string) map[string][]string {
	headers := make(map[string][]string)
	if input == "" {
		return headers
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = append(headers[key], value)
		}
	}

	return headers
}
