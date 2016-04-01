package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
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

func serverGetImageFile(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	for k, _ := range params {
		switch k {
		case "account":
			fallthrough
		case "channel":
			sendResponse(w, InsufficientServerVersion, map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			})
			return

		default:
			sendResponse(w, InvalidParameter, map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter: %s", k),
			})
			return
		}
	}

	filename, exists := getImageFile(path)
	if exists {
		serveFile(w, r, filename, "application/octet-stream")
	} else {
		sendResponse(w, ResourceNotFound, map[string]interface{}{
			"code":    "ResourceNotFound",
			"message": "No such image",
		})
	}
}
