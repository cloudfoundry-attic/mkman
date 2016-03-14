package validators

import (
	"fmt"
	"os"
)

type directory struct {
}

func Directory() Validator {
	return &directory{}
}

func (a *directory) ComposableName() string {
	return "path to directory"
}

func (a *directory) Validate(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	fileInfo, err := os.Stat(convertedObject)
	if os.IsNotExist(err) || !fileInfo.Mode().IsDir() {
		return fmt.Errorf("value must be %s: '%s'", a.ComposableName(), vt.object)
	}

	return nil
}
