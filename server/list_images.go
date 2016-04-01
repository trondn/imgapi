package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/trondn/imgapi/errorcodes"
	"github.com/trondn/imgapi/manifest"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func doListImages(path string, w http.ResponseWriter, r *http.Request) (int, map[string]interface{}) {
	parameters, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		message := map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to parse query parameters: %v", err),
		}
		return errorcodes.InternalError, message
	}

	keys := []string{
		"account",
		"channel",
		"owner",
		"state",
		"name",
		"version",
		"public",
		"os",
		"type",
		"billing_tag",
		"limit",
		"marker",
	}

	for k, _ := range parameters {
		if !stringInSlice(k, keys) {
			message := map[string]interface{}{
				"code":    "InternalError",
				"message": fmt.Sprintf("Invalid key \"%s\"", k),
			}
			return errorcodes.InternalError, message
		}
	}

	// Build up the filter, iterate the spool and generate the restult

	var buffer bytes.Buffer
	buffer.WriteString("[")

	first := true
	dir, _ := ioutil.ReadDir(path)
	for i := 0; i < len(dir); i++ {
		fileinfo := dir[i]
		if !fileinfo.IsDir() {
			log.Printf("Skipping %s (not a directory)", fileinfo.Name())
			continue
		}

		manifest, err := manifest.Load(path + "/" + fileinfo.Name() + "/manifest.json")
		if err != nil {
			log.Printf("Failed to load manifest %s: %e", manifest, err)
			continue
		}

		state, ok := manifest["state"]
		if !ok {
			log.Printf("No state in manifest for: %s", path+"/"+fileinfo.Name()+"/manifest.json")
			continue
		}

		include := false

		// @todo add filter!!
		if state == "active" {
			include = true
		}

		// @TODO check if it match the filter, for now lets assume no filter is provided
		if include {
			if first {
				first = false
			} else {
				buffer.WriteString(",")
			}

			a, _ := json.MarshalIndent(manifest, "  ", "  ")
			buffer.Write(a)
		}
	}
	buffer.WriteString("]")

	var dat map[string]interface{}
	json.Unmarshal(buffer.Bytes(), &dat)
	return errorcodes.Success, dat

}

func ListImages(path string, w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")
	h.Set("Content-Type", "application/json; charset=utf-8")

	code, m := doListImages(path, w, r)
	w.WriteHeader(code)
	a, _ := json.MarshalIndent(m, "", "  ")
	w.Write(a)
}
