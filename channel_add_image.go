package main

import (
	"net/http"
)

func serverChannelAddImage(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, ResourceNotFound, map[string]interface{}{
		"code":    "ResourceNotFound",
		"message": "No support for adding images to channels",
	})
}
