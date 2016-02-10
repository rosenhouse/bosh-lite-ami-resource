package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/rosenhouse/bosh-lite-ami-resource/lib"
)

func main() {
	// dstDir := os.Args[1]

	stdinBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var inData struct {
		Source  lib.Source
		Version struct{}
		Params  struct{}
	}

	err = json.Unmarshal(stdinBytes, &inData)
	if err != nil {
		panic(err)
	}

	var outData struct {
		Version  *struct{} `json:"version,omitempty"`
		Metadata []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"metadata,omitempty"`
	}

	stdoutBytes, err := json.Marshal(outData)
	if err != nil {
		panic(err)
	}

	_, err = os.Stdout.Write(stdoutBytes)
	if err != nil {
		panic(err)
	}
}
