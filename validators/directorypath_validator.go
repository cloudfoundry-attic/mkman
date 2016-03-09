package validators

import (
	"fmt"
	"os"
)

type directoryValidator struct {
}

func NewDirectoryValidator() Validator {
	return &directoryValidator{}
}

func (a *directoryValidator) Name() string {
	return "path to directory"
}

func (a *directoryValidator) Validate(vt ValidationTarget) error {
	fileInfo, err := os.Stat(vt.object)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: '%s'", "directory", vt.object)
		// Return when next error does not make sense
	}

	if err != nil {
		panic(err)
	}

	if !fileInfo.Mode().IsDir() {
		return fmt.Errorf("value must be path to %s: '%s'", "directory", vt.object)
	}

	return nil
}
