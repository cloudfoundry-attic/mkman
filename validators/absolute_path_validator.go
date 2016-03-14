package validators

import (
	"fmt"
	"path/filepath"
)

type absolutePathValidator struct {
}

func AbsolutePathValidator() Validator {
	return &absolutePathValidator{}
}

func (a *absolutePathValidator) ComposableName() string {
	return "absolute path"
}

func (a *absolutePathValidator) Validate(vt ValidationTarget) error {
	err := validateIsAbsPath(vt)
	if err != nil {
		return err
	}
	return nil
}

func validateIsAbsPath(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	if filepath.IsAbs(convertedObject) {
		return nil
	}
	return fmt.Errorf("value must be absolute path: '%s'", vt.object)
}
