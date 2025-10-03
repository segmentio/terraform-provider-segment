package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/segmentio/terraform-provider-segment/internal/provider/models"
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

// mergeSettings merges config settings with remote settings, preserving only the keys defined in config.
func mergeSettings(configSettings, remoteSettings jsontypes.Normalized, isWarehouse bool) (jsontypes.Normalized, error) {
	var configMap map[string]interface{}
	if diags := configSettings.Unmarshal(&configMap); diags.HasError() {
		return jsontypes.NewNormalizedNull(), fmt.Errorf("failed to unmarshal config settings: %s", diags.Errors())
	}

	var remoteMap map[string]interface{}
	if diags := remoteSettings.Unmarshal(&remoteMap); diags.HasError() {
		return jsontypes.NewNormalizedNull(), fmt.Errorf("failed to unmarshal remote settings: %s", diags.Errors())
	}

	// Create merged map with only config-defined keys that exist in remote (to detect drift)
	// Keys in config but not in remote are excluded (they don't exist or aren't supported)
	merged := make(map[string]interface{})
	for key := range configMap {
		if isWarehouse && key == "password" { // Warehouses do not output password in the response
			merged[key] = configMap[key]
		} else if value, exists := remoteMap[key]; exists {
			strValue, ok := value.(string)
			if ok && strings.Contains(strValue, "•") { // If the secret is censored, do not update it
				merged[key] = configMap[key]
			} else {
				merged[key] = value
			}
		}
	}

	result, err := models.GetSettings(merged)
	if err != nil {
		return jsontypes.Normalized{}, fmt.Errorf("failed to merge settings: %s", err.Error())
	}

	return result, nil
}
