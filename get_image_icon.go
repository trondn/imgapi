package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func doServerGetImageIcon(path string, params url.Values) (int, map[string]interface{}) {
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

	m, err := LoadManifest(path + "/manifest.json")
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return InternalError, message
	}

	icon, ok := m["icon"]
	if !ok || icon == false {
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "Image does not have an icon",
		}
		return ResourceNotFound, message
	}

	filename, _ := getIconFile(path)
	if len(filename) == 0 {
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "No such image",
		}
		return ResourceNotFound, message
	}

	return Success, nil
}

func serverGetImageIcon(w http.ResponseWriter, r *http.Request, params url.Values, path string) {

	code, content := doServerGetImageIcon(path, params)
	if code == Success {
		filename, content_type := getIconFile(path)
		serveFile(w, r, filename, content_type)
	} else {
		sendResponse(w, code, content)
	}
}
