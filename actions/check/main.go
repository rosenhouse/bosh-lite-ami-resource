package main

import (
	"log"
	"os"

	"github.com/rosenhouse/bosh-lite-ami-resource/lib"
)

func main() {
	jsonIO := lib.JsonIO{
		InStream:  os.Stdin,
		OutStream: os.Stdout,
	}

	var inData struct {
		Source  lib.Source
		Version lib.Version
	}

	if err := jsonIO.ReadJSON(&inData); err != nil {
		log.Fatalf("%s", err)
	}

	resource := lib.NewResource(inData.Source)
	versions, err := resource.Check(inData.Version)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := jsonIO.WriteJSON(versions); err != nil {
		log.Fatalf("%s", err)
	}
}
