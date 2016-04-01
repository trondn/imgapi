package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/errorcodes"
	"github.com/trondn/imgapi/manifest"
)

func getIconFile(path string) (filename string, content_type string) {
	// @todo loop this :-)
	filename = path + "/icon.png"
	_, err := os.Stat(filename)
	if err == nil {
		return filename, "image/png"
	}

	filename = path + "/icon.jpg"
	_, err = os.Stat(filename)
	if err == nil {
		return filename, "image/jpg"
	}

	filename = path + "/icon.gif"
	_, err = os.Stat(filename)
	if err == nil {
		return filename, "image/gif"
	}

	return "", ""
}

func doDeleteImageIcon(path string, params url.Values) (int, map[string]interface{}) {
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

	filename, _ := getIconFile(path)
	if len(filename) == 0 {
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "No such image",
		}
		return errorcodes.ResourceNotFound, message
	}

	m, err := manifest.Load(path + "/manifest.json")
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return errorcodes.InternalError, message
	}
	m["icon"] = false
	err = manifest.Store(path+"/manifest.json", m)
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to store manifest: %v", err),
		}
		return errorcodes.InternalError, message
	}

	os.Remove(filename)
	return errorcodes.Success, m
}

func DeleteImageIcon(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")
	h.Set("Content-Type", "application/json; charset=utf-8")

	code, m := doDeleteImageIcon(path, params)

	w.WriteHeader(code)
	a, _ := json.MarshalIndent(m, "", "  ")
	w.Write(a)
}
