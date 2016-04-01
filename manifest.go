package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func LoadManifest(path string) (manifest map[string]interface{}, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return manifest, err
	}

	err = json.Unmarshal(content, &manifest)
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}

func StoreManifest(path string, manifest map[string]interface{}) (err error) {
	content, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ManifestValidateType(value interface{}) error {
	// value should be string!!
	switch value.(type) {
	case string:
		break
	default:
		return errors.New("Invalid type for \"type\"")
	}

	legal := []string{"zone-dataset", "lx-dataset", "zvol", "other"}
	if !stringInSlice(value.(string), legal) {
		return errors.New(fmt.Sprintf("Invalid value specified for \"type\": \"%v\"", value))
	}

	return nil
}

func ManifestValidateOs(value interface{}) error {
	// value should be string!!
	switch value.(type) {
	case string:
		break
	default:
		return errors.New("Invalid type for \"os\"")
	}

	legal := []string{"smartos", "windows", "linux", "bsd"}
	if !stringInSlice(value.(string), legal) {
		return errors.New(fmt.Sprintf("Invalid value specified for \"os\": \"%v\"", value))
	}

	return nil
}

func ManifestValidateCompression(value interface{}) error {
	// value should be string!!
	switch value.(type) {
	case string:
		break
	default:
		return errors.New("Invalid type for \"compression\"")
	}

	legal := []string{"bzip2", "gzip", "none"}
	if !stringInSlice(value.(string), legal) {
		return errors.New(fmt.Sprintf("Invalid value specified for \"compression\": \"%v\"", value))
	}

	return nil
}
