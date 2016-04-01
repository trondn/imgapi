package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func doServerDisableImage(path string, params url.Values) (int, map[string]interface{}) {
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

	m["disabled"] = true
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

func serverDisableImage(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	code, content := doServerDisableImage(path, params)
	sendResponse(w, code, content)

}
