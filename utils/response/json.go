package response

import (
	"encoding/json"
	"fmt"
)

func FormatJSONResponse(rawJSON string) string {
	var jsonResponse map[string]interface{}

	// Unmarshal the raw JSON string into a Go map
	if err := json.Unmarshal([]byte(rawJSON), &jsonResponse); err != nil {
		return fmt.Sprintf("Error decoding JSON: %v", err)
	}

	// Marshal the JSON with indentation
	prettyJSON, err := json.MarshalIndent(jsonResponse, "\t", "\t") // Using tabs for indentation
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}

	// Return the formatted JSON string
	return string(prettyJSON)
}
