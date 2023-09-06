package provider

import (
	"bytes"
	"encoding/json"
	"io"
)

func getError(err error, body io.ReadCloser) string {
	parsedBody, readErr := io.ReadAll(body)
	if readErr != nil {
		return err.Error()
	}
	var formattedBody bytes.Buffer
	jsonErr := json.Indent(&formattedBody, parsedBody, "", "  ")
	if jsonErr != nil {
		return err.Error()
	}
	return err.Error() + "\n" + formattedBody.String()
}
