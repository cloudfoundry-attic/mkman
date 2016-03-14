package validators

import (
	"fmt"
	"os"
)

type existenceValidator struct {
}

func ExistenceValidator() Validator {
	return &existenceValidator{}
}

func (ev *existenceValidator) ComposableName() string {
	return "present on filesystem"
}

func (ev *existenceValidator) Validate(vt ValidationTarget) error {
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
