package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/errorcodes"
)

func getImageFile(path string) (filename string, ok bool) {
	ext := []string{".bz2", ".gz", ""}

	for i := 0; i < len(ext); i++ {
		filename = path + "/image" + ext[i]
		_, err := os.Stat(filename)
		if err == nil {
			return filename, true
		}
	}

	return filename, false
}

func GetImageFile(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")

	for k, _ := range params {
		switch k {
		case "account":
		case "channel":
			message := map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			}

			h.Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(errorcodes.InsufficientServerVersion)
			a, _ := json.MarshalIndent(message, "", "  ")
			w.Write(a)
			return

		default:
			message := map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter: %s", k),
			}
			h.Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(errorcodes.InvalidParameter)
			a, _ := json.MarshalIndent(message, "", "  ")
			w.Write(a)
			return
		}
	}

	filename, exists := getImageFile(path)
	if exists {
		http.ServeFile(w, r, filename)
	} else {
		h.Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(errorcodes.ResourceNotFound)
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "No such image",
		}
		a, _ := json.MarshalIndent(message, "", "  ")
		w.Write(a)
	}
}
