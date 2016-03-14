package validators

import (
	"fmt"
	"path/filepath"
)

type absolutePath struct {
}

func AbsolutePath() Validator {
	return &absolutePath{}
}

func (a *absolutePath) ComposableName() string {
	return "absolute path"
}

func (a *absolutePath) Validate(vt ValidationTarget) error {
	err := validateIsAbsPath(vt)
	if err != nil {
		return err
	}
	return nil
}

func validateIsAbsPath(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		// We don't want to return an error here as it would look like a validation
		// failure instead of a coding error.
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	if filepath.IsAbs(convertedObject) {
		return nil
	}
	return fmt.Errorf("value must be absolute path: '%s'", vt.object)
}
