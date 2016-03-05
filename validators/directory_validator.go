package validators

import (
	"fmt"
	"os"
)

type directoryValidator struct {
}

func DirectoryValidator() Validator {
	return &directoryValidator{}
}

func (a *directoryValidator) ComposableName() string {
	return "path to directory"
}

func (a *directoryValidator) Validate(vt ValidationTarget) error {
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
