package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/contrib"
)

func addDefaultValue(key string, value interface{}, manifest map[string]interface{}) {
	_, exists := manifest[key]
	if !exists {
		manifest[key] = value
	}
}

func doServerCreateImage(w http.ResponseWriter, r *http.Request, params url.Values, datadir string) (int, map[string]interface{}) {
	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed to read body: %e", err)
		return InternalError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to read body: %v", err),
		}
	}

	var m map[string]interface{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		log.Printf("Failed to parse payload: %e", err)
		return InternalError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to decode body: %v", err),
		}
	}

	mandatory := []string{
		"name",
		"version",
		"type",
		"os",
	}

	for i := 0; i < len(mandatory); i++ {
		_, present := m[mandatory[i]]
		if !present {
			return InvalidParameter, map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Mandatory key \"%s\" is not present", mandatory[i]),
			}
		}
	}

	// Ok, lets walk through the parameters given
	for k, v := range m {
		switch k {
		case "owner":
			fallthrough
		case "name":
			fallthrough
		case "version":
			fallthrough
		case "description":
			fallthrough
		case "homepage":
			fallthrough
		case "eula":
			fallthrough
		case "disabled":
			fallthrough
		case "public":
			break

		case "type":
			err = ManifestValidateType(v)
			break

		case "os":
			err = ManifestValidateOs(v)
			break

		case "origin":
			fallthrough
		case "acl":
			fallthrough
		case "requirements":
			fallthrough
		case "users":
			fallthrough
		case "billing_tags":
			fallthrough
		case "traits":
			fallthrough
		case "tags":
			fallthrough
		case "generate_passwords":
			fallthrough
		case "inherited_directories":
			fallthrough
		case "nic_driver":
			fallthrough
		case "disk_driver":
			fallthrough
		case "cpu_type":
			fallthrough
		case "image_size":
			break

		default:
			err = errors.New(fmt.Sprintf("Unknown parameter: %s", k))
		}

		if err != nil {
			return InvalidParameter, map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("%v", err),
			}
		}
	}

	uuid, _ := contrib.NewUUID()
	addDefaultValue("uuid", uuid, m)
	addDefaultValue("state", "unactivated", m)
	addDefaultValue("disabled", false, m)
	addDefaultValue("public", false, m)
	addDefaultValue("v", 2, m)

	// Validate that the uuid don't exists
	path := datadir + "/" + uuid

	err = os.Mkdir(path, 0777)
	if err != nil {
		if os.IsExist(err) {
			return ImageUuidAlreadyExists, map[string]interface{}{
				"code":    "ImageUuidAlreadyExists",
				"message": "Uuid already exists",
			}
		}

		return InternalError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Internal error: %v", err),
		}
	}

	err = StoreManifest(path+"/manifest.json", m)
	if err != nil {
		_ = os.RemoveAll(path)
		return InternalError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to write manifest: %v", err),
		}
	}

	return Success, m
}

func serverCreateImage(w http.ResponseWriter, r *http.Request, params url.Values, datadir string) {
	code, content := doServerCreateImage(w, r, params, datadir)
	sendResponse(w, code, content)
}
