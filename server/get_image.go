package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/trondn/imgapi/errorcodes"
	"github.com/trondn/imgapi/manifest"
)

func doGetImage(path string, params url.Values) (int, map[string]interface{}) {
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

	m, err := manifest.Load(path + "/manifest.json")
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return errorcodes.InternalError, message
	}

	return errorcodes.Success, m
}

func GetImage(w http.ResponseWriter, r *http.Request, params url.Values, path string) {

	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")
	h.Set("Content-Type", "application/json; charset=utf-8")

	code, m := doGetImage(path, params)

	w.WriteHeader(code)
	a, _ := json.MarshalIndent(m, "", "  ")
	w.Write(a)
}
