package validators

import (
	"fmt"
	"os"
)

type existsOnFilesystem struct {
}

func ExistsOnFilesystem() Validator {
	return &existsOnFilesystem{}
}

func (ev *existsOnFilesystem) ComposableName() string {
	return "present on filesystem"
}

func (ev *existsOnFilesystem) Validate(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	_, err := os.Stat(convertedObject)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s must be %s: '%s'", vt.name, ev.ComposableName(), vt.object)
	}
	return nil
}
