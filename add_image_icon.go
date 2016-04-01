package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func doServerAddImageIcon(path string, params url.Values, header http.Header, reader io.Reader) (int, map[string]interface{}) {
	content_type := header.Get("Content-Type")
	var extension string

	switch content_type {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "image/gif":
		extension = ".gif"
	case "":
		message := map[string]interface{}{
			"code":    "InvalidParameter",
			"message": "Content-Type not present",
		}
		return InvalidParameter, message
	default:
		message := map[string]interface{}{
			"code":    "InvalidParameter",
			"message": fmt.Sprintf("Unknown Content-Type \"%s\"", content_type),
		}
		return InvalidParameter, message
	}

	var expectedsha1 string
	for k, v := range params {
		switch k {
		case "account":
			fallthrough
		case "channel":
			fallthrough
		case "storage":
			message := map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			}
			return InsufficientServerVersion, message

		case "sha1":
			expectedsha1 = v[0]

		default:
			message := map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter: %s", k),
			}
			return InvalidParameter, message
		}
	}

	filename := path + "/icon" + extension
	writer, err := os.Create(filename)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to create image file: %v", err),
		}
		return InternalError, message
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to store image file: %v", err),
		}
		return InternalError, message
	}

	if len(expectedsha1) > 0 {
		sha1, err := GetSha1Sum(filename)
		if err != nil {
			os.Remove(filename)
			message := map[string]interface{}{
				"code":    "InternalError",
				"message": fmt.Sprintf("Failed to generate sha1: %v", err),
			}
			return InternalError, message
		}

		if expectedsha1 != sha1 {
			os.Remove(filename)
			message := map[string]interface{}{
				"code":    "InternalError",
				"message": fmt.Sprintf("Incorrect SHA. expected \"%s\" got \"%s\"", expectedsha1, sha1),
			}
			return InternalError, message
		}
	}

	m, err := LoadManifest(path + "/manifest.json")
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return InternalError, message
	}
	m["icon"] = true
	err = StoreManifest(path+"/manifest.json", m)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to store manifest: %v", err),
		}
		return InternalError, message
	}

	return Success, m
}

func serverAddImageIcon(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	code, content := doServerAddImageIcon(path, params, r.Header, r.Body)
	sendResponse(w, code, content)
}
