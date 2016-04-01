package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/errorcodes"
	"github.com/trondn/imgapi/manifest"
)

func doAddImageFile(path string, params url.Values, reader io.Reader) (int, map[string]interface{}) {
	var expectedsha1 string
	var compression string
	for k, v := range params {
		switch k {
		case "account":
		case "channel":
		case "storage":
		case "dataset_guid":
			message := map[string]interface{}{
				"code":    "InsufficientServerVersion",
				"message": "The server does not support \"account\" and \"channel\"",
			}
			return errorcodes.InsufficientServerVersion, message

		case "compression":
			compression = v[0]
			switch compression {
			case "gzip":
			case "bzip2":
				break
			default:
				message := map[string]interface{}{
					"code":    "InvalidParameter",
					"message": "compression may be gzip of bzip2",
				}
				return errorcodes.InvalidParameter, message
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
			return errorcodes.InvalidParameter, message
		}
	}

	manifestfile := path + "/manifest.json"

	m, err := manifest.Load(manifestfile)
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to load manifest: %v", err),
		}
		return errorcodes.InternalError, message
	}

	if m["state"] == "active" {
		message := map[string]interface{}{
			"code":    "ImageAlreadyActivated",
			"message": "Can't replace file for an active image",
		}
		return errorcodes.ImageAlreadyActivated, message
	}

	filename := path + "/image.gz"
	f, err := os.Create(filename)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to create image file: %v", err),
		}
		return errorcodes.InternalError, message
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
			return errorcodes.InternalError, message
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
		return errorcodes.InternalError, message
	}

	stat, err := os.Stat(filename)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to lookup image file: %v", err),
		}
		return errorcodes.InternalError, message
	}

	// ok, generate the SHA1
	sha1sum, err := getSha1Sum(filename)
	if err != nil {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to get SHA1 for image file: %v", err),
		}
		return errorcodes.InternalError, message
	}

	if len(expectedsha1) > 0 && sha1sum != expectedsha1 {
		os.Remove(filename)
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Incorrect SHA. expected \"%s\" got \"%s\"", expectedsha1, sha1sum),
		}
		return errorcodes.InternalError, message
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
	err = manifest.Store(manifestfile, m)
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

func AddImageFile(w http.ResponseWriter, r *http.Request, params url.Values, path string) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")
	h.Set("Content-Type", "application/json; charset=utf-8")

	code, m := doAddImageFile(path, params, r.Body)

	w.WriteHeader(code)
	a, _ := json.MarshalIndent(m, "", "  ")
	w.Write(a)
}
