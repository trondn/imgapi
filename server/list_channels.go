package server

import (
	"encoding/json"
	"net/http"

	"github.com/trondn/imgapi/errorcodes"
)

func ListChannels(w http.ResponseWriter, r *http.Request) {
	payload := map[string]interface{}{
		"code":    "ResourceNotFound",
		"message": "/channels does not exist",
	}

	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errorcodes.ResourceNotFound)

	a, _ := json.MarshalIndent(payload, "", "  ")
	w.Write(a)
}
