package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
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

func doServerDeleteImageIcon(path string, params url.Values) (int, map[string]interface{}) {
	for k, _ := range params {
		switch k {
		case "account":
			fallthrough
		case "channel":
			message := map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			}
			return InsufficientServerVersion, message

		default:
			message := map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter: %s", k),
			}
			return InvalidParameter, message
		}
	}

	filename, _ := getIconFile(path)
	if len(filename) == 0 {
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "No such image",
		}
		return ResourceNotFound, message
	}

	m, err := LoadManifest(path + "/manifest.json")
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return InternalError, message
	}
	m["icon"] = false
	err = StoreManifest(path+"/manifest.json", m)
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to store manifest: %v", err),
		}
		return InternalError, message
	}

	os.Remove(filename)
	return Success, m
}

func serverDeleteImageIcon(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	code, content := doServerDeleteImageIcon(path, params)
	sendResponse(w, code, content)
}
