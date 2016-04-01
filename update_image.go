package main

import (
	"net/http"
	"net/url"
)

func serverUpdateImage(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	sendResponse(w, InsufficientServerVersion, map[string]interface{}{
		"code":    "InsufficientServerVersion",
		"message": "Not implemented yet",
	})
}
