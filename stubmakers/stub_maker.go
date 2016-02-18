package stubmakers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//go:generate counterfeiter . StubMaker
type StubMaker interface {
	MakeStub() (string, error)
}

func marshalTempStub(objectToMarshal interface{}, filename string) (string, error) {
	stubContents, err := yaml.Marshal(objectToMarshal)
	if err != nil {
		// We cannot test this because it is too hard to get the marshaller to
		// return an error
		return "", nil
	}

	intermediateDir, err := ioutil.TempDir("", "")
	if err != nil {
		// We cannot test this because it is too hard to get TempDir to return error
		return "", err
	}

	stubPath := filepath.Join(intermediateDir, filename)
	err = ioutil.WriteFile(stubPath, stubContents, os.ModePerm)
	if err != nil {
		// We cannot test this because it is hard to simulate an error with
		// WriteFile
		return "", err
	}

	return stubPath, nil
}
