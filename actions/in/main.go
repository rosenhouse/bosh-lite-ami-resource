package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/rosenhouse/bosh-lite-ami-resource/lib"
)

func checkDirExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("exists but not directory: %s", path)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("expected exactly 1 arg")
	}

	dstDir := os.Args[1]
	if err := checkDirExists(dstDir); err != nil {
		log.Fatalf("%s", err)
	}

	jsonIO := lib.JsonIO{
		InStream:  os.Stdin,
		OutStream: os.Stdout,
	}

	var inData struct {
		Source  lib.Source
		Version lib.Version
		Params  struct{}
	}

	if err := jsonIO.ReadJSON(&inData); err != nil {
		log.Fatalf("%s", err)
	}

	resource := lib.NewResource(inData.Source)
	ami, err := resource.In(inData.Version)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := ioutil.WriteFile(
		filepath.Join(dstDir, "ami"),
		[]byte(ami), 0644); err != nil {
		log.Fatalf("%s", err)
	}

	var outData struct {
		Version  lib.Version    `json:"version,omitempty"`
		Metadata []lib.Metadata `json:"metadata,omitempty"`
	}

	outData.Version = inData.Version
	outData.Metadata = []lib.Metadata{
		{Name: "box_version", Value: outData.Version.BoxVersion},
		{Name: "ami", Value: ami},
		{Name: "region", Value: resource.SourceConfig.Region},
		{Name: "box_name", Value: resource.SourceConfig.BoxName},
	}

	if err := jsonIO.WriteJSON(outData); err != nil {
		log.Fatalf("%s", err)
	}
}
