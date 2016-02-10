package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/rosenhouse/bosh-lite-ami-resource/lib"
)

func main() {
	stdinBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var inData struct {
		Source  lib.Source
		Version lib.Version
	}

	err = json.Unmarshal(stdinBytes, &inData)
	if err != nil {
		log.Fatalf("%s", err)
	}

	resource := lib.NewResource(inData.Source)
	versions, err := resource.Check(inData.Version)
	if err != nil {
		log.Fatalf("%s", err)
	}

	stdoutBytes, err := json.Marshal(versions)
	if err != nil {
		log.Fatalf("%s", err)
	}

	_, err = os.Stdout.Write(stdoutBytes)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
