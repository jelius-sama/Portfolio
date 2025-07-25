package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// TODO: Work needed
func PurgeCache(urls []string) error {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	payload := map[string]interface{}{
		"files": urls,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/purge_cache", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to purge cache: %s", resp.Status)
	}

	return nil
}
