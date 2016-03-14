package validators

import (
	"fmt"
	"os"
)

type file struct {
}

func File() FileTypeValidator {
	return &file{}
}

func (a *file) ComposableName() string {
	return "a file"
}

func (a *file) Validate(vt ValidationTarget) error {
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

func (a *file) FileType() string {
	return a.ComposableName()
}
