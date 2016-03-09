package validators

import (
	"fmt"
	"os"
)

type filepathValidator struct {
}

func NewFilepathValidator() Validator {
	return &filepathValidator{}
}

func (a *filepathValidator) Name() string {
	return "path to file"
}

func (a *filepathValidator) Validate(vt ValidationTarget) error {
	fileInfo, err := os.Stat(vt.object)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: '%s'", "file", vt.object)
		// Return when next error does not make sense
	}

	if err != nil {
		panic(err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("value must be path to %s: '%s'", "file", vt.object)
	}

	return nil
}
