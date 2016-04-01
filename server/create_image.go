package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/trondn/imgapi/common"
	"github.com/trondn/imgapi/errorcodes"
	"github.com/trondn/imgapi/manifest"
)

func addDefaultValue(key string, value interface{}, manifest map[string]interface{}) {
	_, exists := manifest[key]
	if !exists {
		manifest[key] = value
	}
}

func doCreateImage(w http.ResponseWriter, r *http.Request, params url.Values, datadir string) (int, map[string]interface{}) {
	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed to read body: %e", err)
		return errorcodes.InternalError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to read body: %v", err),
		}
	}

	var m map[string]interface{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		log.Printf("Failed to parse payload: %e", err)
		return errorcodes.InternalError, map[string]interface{}{
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
			return errorcodes.InvalidParameter, map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("Mandatory key \"%s\" is not present", mandatory[i]),
			}
		}
	}

	// Ok, lets walk through the parameters given
	for k, v := range m {
		switch k {
		case "owner":
		case "name":
		case "version":
		case "description":
		case "homepage":
		case "eula":
		case "disabled":
		case "public":
			break

		case "type":
			err = manifest.ValidateType(v)
			break

		case "os":
			err = manifest.ValidateOs(v)
			break

		case "origin":
		case "acl":
		case "requirements":
		case "users":
		case "billing_tags":
		case "traits":
		case "tags":
		case "generate_passwords":
		case "inherited_directories":
		case "nic_driver":
		case "disk_driver":
		case "cpu_type":
		case "image_size":
			break

		default:
			err = errors.New(fmt.Sprintf("Unknown parameter: %s", k))
		}

		if err != nil {
			return errorcodes.InvalidParameter, map[string]interface{}{
				"code":    "InvalidParameter",
				"message": fmt.Sprintf("%v", err),
			}
		}
	}

	uuid, _ := common.NewUUID()
	addDefaultValue("uuid", uuid, m)
	addDefaultValue("state", "unactivated", m)
	addDefaultValue("disabled", false, m)
	addDefaultValue("public", false, m)

	// Validate that the uuid don't exists
	path := datadir + "/" + uuid

	err = os.Mkdir(path, 0777)
	if err != nil {
		if os.IsExist(err) {
			return errorcodes.ImageUuidAlreadyExists, map[string]interface{}{
				"code":    "ImageUuidAlreadyExists",
				"message": "Uuid already exists",
			}
		}

		return errorcodes.InternalError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Internal error: %v", err),
		}
	}

	err = manifest.Store(path+"/manifest.json", m)
	if err != nil {
		_ = os.RemoveAll(path)
		return errorcodes.InternalError, map[string]interface{}{
			"code":    "InternalError",
			"message": fmt.Sprintf("Failed to write manifest: %v", err),
		}
	}

	return errorcodes.Success, m
}

func CreateImage(w http.ResponseWriter, r *http.Request, params url.Values, datadir string) {
	h := w.Header()
	h.Set("Server", "Norbye Public Images Repo")
	h.Set("Content-Type", "application/json; charset=utf-8")

	code, m := doCreateImage(w, r, params, datadir)

	w.WriteHeader(code)
	a, _ := json.MarshalIndent(m, "", "  ")
	w.Write(a)
}
