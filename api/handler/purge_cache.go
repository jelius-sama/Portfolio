package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"KazuFolio/util"
	"strings"
)

// TODO: Test needed
func PurgeCache(w http.ResponseWriter, r *http.Request) {
	VerifySudo(w, r)

	var body struct {
		Paths []string `json:"paths"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Paths) == 0 {
		util.WriteError(w, http.StatusBadRequest, "400 Bad Request", util.AddrOf("invalid request"))
		return
	}

	base := os.Getenv("host")

	for i, path := range body.Paths {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		body.Paths[i] = base + path
	}

	payload, _ := json.Marshal(map[string][]string{"files": body.Paths})
	req, _ := http.NewRequest("POST",
		"https://api.cloudflare.com/client/v4/zones/"+os.Getenv("CF_ZONE_ID")+"/purge_cache",
		bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CF_API_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("cache purge failed"))
		return
	}

	util.WriteSuccess(w, http.StatusOK, "200 OK", util.AddrOf("cache purged successfully"))
}
