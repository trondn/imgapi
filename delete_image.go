package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func doServerDeleteImage(path string, params url.Values) (int, map[string]interface{}) {
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

	_, err := os.Stat(path)
	if err != nil {
		message := map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "The image does not exist",
		}
		return ResourceNotFound, message
	}

	os.RemoveAll(path)
	return NoContent, nil
}

func serverDeleteImage(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	code, content := doServerDeleteImage(path, params)
	sendResponse(w, code, content)
}
