package parser

import (
	"encoding/json"
	"fmt"
	"net/http"

	"KazuFolio/api/handler"
	"KazuFolio/util"
)

func SSRAchievementsPage(fullPath string) (string, error, int) {
	achievementsStats := handler.GetAchievementsStatsInternal()

	resp_body, err := json.Marshal(achievementsStats)
	if err != nil {
		return "", fmt.Errorf("failed to read API response: %w", err), http.StatusInternalServerError
	}

	metaHTML, metaJSON, err := ParseStaticMetadataForPaths([]string{"*", "/achievements"})

	if err != nil {
		return "", fmt.Errorf("failed to parse metadata: %w", err), http.StatusInternalServerError
	}

	// Step 7: Prepare final JSON for hydration script
	scriptPayload := map[string]any{
		"path":     fullPath,
		"metadata": metaJSON,
		"api_resp": json.RawMessage(resp_body),
	}
	scriptJSON, err := json.Marshal(scriptPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal script payload: %w", err), http.StatusInternalServerError
	}

	scriptHTML := `<script id="__SERVER_DATA__" type="application/json">` + util.HTMLEscape(scriptJSON) + `</script>`

	// Step 8: Return combined HTML
	return scriptHTML + "\n" + metaHTML, nil, http.StatusOK
}
