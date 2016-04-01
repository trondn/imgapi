package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/errorcodes"
)

func doDeleteImage(path string, params url.Values) (int, map[string]interface{}) {
	for k, _ := range params {
		switch k {
		case "account":
		case "channel":
			message := map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			}
			return errorcodes.InsufficientServerVersion, message

		default:
			message := map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter: %s", k),
			}
			return errorcodes.InvalidParameter, message
		}
	}

	_, err := os.Stat(path)
	if err != nil {
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "The image does not exist",
		}
		return errorcodes.ResourceNotFound, message
	}

	os.RemoveAll(path)
	return errorcodes.NoContent, nil
}

func DeleteImage(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")

	code, m := doDeleteImage(path, params)

	w.WriteHeader(code)
	if m != nil {
		h.Set("Content-Type", "application/json; charset=utf-8")
		a, _ := json.MarshalIndent(m, "", "  ")
		w.Write(a)
	}
}
