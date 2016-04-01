package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/errorcodes"
	"github.com/trondn/imgapi/manifest"
)

func doAddImageIcon(path string, params url.Values, header http.Header, reader io.Reader) (int, map[string]interface{}) {
	content_type := header.Get("Content-Type")
	var extension string

	switch content_type {
	case "image/jpeg":
		extension = ".jpg"
		break
	case "image/png":
		extension = ".png"
		break
	case "image/gif":
		extension = ".gif"
		break
	case "":
		message := map[string]interface{}{
			"code":    "InvalidParameter",
			"message": "Content-Type not present",
		}
		return errorcodes.InvalidParameter, message
	default:
		message := map[string]interface{}{
			"code":    "InvalidParameter",
			"message": fmt.Sprintf("Unknown Content-Type \"%s\"", content_type),
		}
		return errorcodes.InvalidParameter, message
	}

	var expectedsha1 string
	for k, v := range params {
		switch k {
		case "account":
		case "channel":
		case "storage":
			message := map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			}
			return errorcodes.InsufficientServerVersion, message

		case "sha1":
			expectedsha1 = v[0]
			break

		default:
			message := map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter: %s", k),
			}
			return errorcodes.InvalidParameter, message
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
		return errorcodes.InternalError, message
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to store image file: %v", err),
		}
		return errorcodes.InternalError, message
	}

	if len(expectedsha1) > 0 {
		sha1, err := getSha1Sum(filename)
		if err != nil {
			os.Remove(filename)
			message := map[string]interface{}{
				"code":    "InternalError",
				"message": fmt.Sprintf("Failed to generate sha1: %v", err),
			}
			return errorcodes.InternalError, message
		}

		if expectedsha1 != sha1 {
			os.Remove(filename)
			message := map[string]interface{}{
				"code":    "InternalError",
				"message": fmt.Sprintf("Incorrect SHA. expected \"%s\" got \"%s\"", expectedsha1, sha1),
			}
			return errorcodes.InternalError, message
		}
	}

	m, err := manifest.Load(path + "/manifest.json")
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return errorcodes.InternalError, message
	}
	m["icon"] = true
	err = manifest.Store(path+"/manifest.json", m)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to store manifest: %v", err),
		}
		return errorcodes.InternalError, message
	}

	return errorcodes.Success, m
}

func AddImageIcon(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")
	h.Set("Content-Type", "application/json; charset=utf-8")

	code, m := doAddImageIcon(path, params, r.Header, r.Body)

	w.WriteHeader(code)
	a, _ := json.MarshalIndent(m, "", "  ")
	w.Write(a)
}
