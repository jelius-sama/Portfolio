package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"KazuFolio/util"
	"strings"
)

func PurgeCache(w http.ResponseWriter, r *http.Request) {
	if !VerifySudo(w, r) {
		return
	}

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

	payload, err := json.Marshal(map[string][]string{"files": body.Paths})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("failed to marshal payload"))
		return
	}

	req, err := http.NewRequest("POST",
		"https://api.cloudflare.com/client/v4/zones/"+os.Getenv("CF_ZONE_ID")+"/purge_cache",
		bytes.NewBuffer(payload))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("failed to create request"))
		return
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("CF_API_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("failed to send request"))
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("cache purge failed"))
		return
	}

	util.WriteSuccess(w, http.StatusOK, "200 OK", util.AddrOf("cache purged successfully"))
}

func PurgeAllCache(w http.ResponseWriter, r *http.Request) {
	if !VerifySudo(w, r) {
		return
	}

	payload, err := json.Marshal(map[string]bool{"purge_everything": true})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("failed to marshal payload"))
		return
	}

	req, err := http.NewRequest("POST",
		"https://api.cloudflare.com/client/v4/zones/"+os.Getenv("CF_ZONE_ID")+"/purge_cache",
		bytes.NewBuffer(payload))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("failed to create request"))
		return
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("CF_API_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("failed to send request"))
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", util.AddrOf("cache purge failed"))
		return
	}

	util.WriteSuccess(w, http.StatusOK, "200 OK", util.AddrOf("all cache purged successfully"))
}
