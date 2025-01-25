package response

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func FormatJSONResponse(body []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return prettyJSON.String()

}
