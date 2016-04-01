package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func doServerEnableImage(path string, params url.Values) (int, map[string]interface{}) {
	for k, _ := range params {
		switch k {
		case "action":
			break
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

	m, err := LoadManifest(path + "/manifest.json")
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("The server failed to load manifest file: %v", err),
		}
		return InternalError, message
	}

	disabled, ok := m["disabled"]
	if !ok {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": "manifest does not contain \"disabled\"",
		}
		return InternalError, message
		disabled = false
	}

	state, ok := m["state"]
	if !ok {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": "manifest does not contain \"state\"",
		}
		return InternalError, message
	}

	if state == "activated" && !disabled.(bool) {
		message := map[string]interface{}{
			"code":    "ImageAlreadyActivated",
			"message": "Image already activated",
		}
		return InternalError, message
	}

	// Verify that I have the image file
	filename := path + "/image.gz"
	_, err = os.Stat(filename)
	if err != nil {
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "No image file",
		}
		return ResourceNotFound, message
	}

	// Ok enable
	m["disabled"] = false
	m["state"] = "activated"
	err = StoreManifest(path+"/manifest.json", m)
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to store manifest file: %v", err),
		}
		return InternalError, message
	}

	return Success, m
}

func serverEnableImage(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	code, content := doServerEnableImage(path, params)
	sendResponse(w, code, content)
}
