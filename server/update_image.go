package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/trondn/imgapi/errorcodes"
)

func UpdateImage(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	payload := map[string]interface{}{
		"code":    "InsufficientServerVersion",
		"message": "Not implemented yet",
	}

	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errorcodes.InsufficientServerVersion)

	a, _ := json.MarshalIndent(payload, "", "  ")
	w.Write(a)
}
