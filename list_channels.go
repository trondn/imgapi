package main

import (
	"net/http"
)

func serverListChannels(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, ResourceNotFound, map[string]interface{}{
		"code":    "ResourceNotFound",
		"message": "/channels does not exist",
	})
}
