package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type JsonIO struct {
	InStream  io.Reader
	OutStream io.Writer
}

func (j *JsonIO) ReadJSON(output interface{}) error {
	inBytes, err := ioutil.ReadAll(j.InStream)
	if err != nil {
		return fmt.Errorf("unable to read from stream: %s", err)
	}

	err = json.Unmarshal(inBytes, &output)
	if err != nil {
		return fmt.Errorf("unable to unmarshal bytes as JSON into %T: %s", output, err)
	}

	return nil
}

func (j *JsonIO) WriteJSON(input interface{}) error {
	outBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("unable to marshal %T into bytes: %s", input, err)
	}

	_, err = j.OutStream.Write(outBytes)
	if err != nil {
		return fmt.Errorf("unable to write to stream: %s", err)
	}

	return nil
}
