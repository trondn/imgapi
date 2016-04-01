package server

import (
	"encoding/json"
	"net/http"

	"github.com/trondn/imgapi/errorcodes"
)

func ChannelAddImage(w http.ResponseWriter, r *http.Request) {
	payload := map[string]interface{}{
		"code":    "ResourceNotFound",
		"message": "No support for adding images to channels",
	}

	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errorcodes.ResourceNotFound)

	a, _ := json.MarshalIndent(payload, "", "  ")
	w.Write(a)
}
