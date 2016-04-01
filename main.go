package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var configuration Configuration

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Failed to get information about current user: %v",
			err)
	}
	// Set up default values
	configfile := usr.HomeDir + "/.imgapi.json"
	server_mode := false

	flag.BoolVar(&server_mode, "s", false, "Server mode")
	flag.StringVar(&configfile, "c", configfile, "Configuration file")
	flag.Parse()

	if len(flag.Args()) != 0 {
		fmt.Println("Usage: imgapi [arguments]")
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Printf("Failed to read %s: %v", configfile, err)
		os.Exit(1)
	}
	err = json.Unmarshal(content, &configuration)
	if err != nil {
		log.Fatalf("Failed to parse JSON: [%s]: %e", content, err)
	}

	if server_mode {
		startImageServer()
	} else {
		log.Fatal("Client API is not implemented")
	}
}
