package validators

import (
	"fmt"
	"os"
)

type fileValidator struct {
}

func FileValidator() Validator {
	return &fileValidator{}
}

func (a *fileValidator) ComposableName() string {
	return "path to file"
}

func (a *fileValidator) Validate(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	fileInfo, err := os.Stat(convertedObject)
	if os.IsNotExist(err) || !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("value must be %s: '%s'", a.ComposableName(), vt.object)
	}

	return nil
}
