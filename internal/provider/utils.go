package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func getError(err error, body *http.Response) string {
	if body == nil {
		return err.Error()
	}
	parsedBody, readErr := io.ReadAll(body.Body)
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
