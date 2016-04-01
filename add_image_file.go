package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func doServerAddImageFile(path string, params url.Values, reader io.Reader) (int, map[string]interface{}) {
	var expectedsha1 string
	var compression string
	for k, v := range params {
		switch k {
		case "account":
			fallthrough
		case "channel":
			fallthrough
		case "storage":
			fallthrough
		case "dataset_guid":
			message := map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			}
			return InsufficientServerVersion, message

		case "compression":
			compression = v[0]
			switch compression {
			case "gzip":
				fallthrough
			case "bzip2":
				break
			default:
				message := map[string]interface{}{
					"code":    "InvalidParameter",
					"message": "compression may be gzip of bzip2",
				}
				return InvalidParameter, message
			}

			break

		case "sha1":
			expectedsha1 = v[0]
			break

		default:
			message := map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Invalid parameter: %s", k),
			}
			return InvalidParameter, message
		}
	}

	manifestfile := path + "/manifest.json"

	m, err := LoadManifest(manifestfile)
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return InternalError, message
	}

	if m["state"] == "active" {
		message := map[string]interface{}{
			"code":    "ImageAlreadyActivated",
			"message": "Can't replace file for an active image",
		}
		return ImageAlreadyActivated, message
	}

	filename := path + "/image.gz"
	f, err := os.Create(filename)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to create image file: %v", err),
		}
		return InternalError, message
	}

	var writer io.Writer
	if len(compression) == 0 {
		writer, err = gzip.NewWriterLevel(f, gzip.BestCompression)
		if err != nil {
			os.Remove(filename)
			message := map[string]interface{}{
				"code":    "InternalError",
				"message": fmt.Sprintf("Failed to create zip stream: %v", err),
			}
			return InternalError, message
		}
		compression = "gzip"
	} else {
		writer = f
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

	stat, err := os.Stat(filename)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to lookup image file: %v", err),
		}
		return InternalError, message
	}

	// ok, generate the SHA1
	sha1sum, err := GetSha1Sum(filename)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to get SHA1 for image file: %v", err),
		}
		return InternalError, message
	}

	if len(expectedsha1) > 0 && sha1sum != expectedsha1 {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Incorrect SHA. expected \"%s\" got \"%s\"", expectedsha1, sha1sum),
		}
		return InternalError, message
	}

	entry := map[string]interface{}{
		"compression": compression,
		"sha1":        sha1sum,
		"size":        stat.Size(),
	}

	files := []map[string]interface{}{
		entry,
	}

	m["files"] = files
	err = StoreManifest(manifestfile, m)
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

func serverAddImageFile(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	code, content := doServerAddImageFile(path, params, r.Body)
	sendResponse(w, code, content)
}
