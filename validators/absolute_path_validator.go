package validators

import (
	"fmt"
	"path/filepath"
)

type absolutePathValidator struct {
}

func NewAbsolutePathValidator() Validator {
	return &absolutePathValidator{}
}

func (a *absolutePathValidator) Name() string {
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
	if filepath.IsAbs(vt.object) {
		return nil
	}
	return fmt.Errorf("value must be absolute path: '%s'", vt.object)
}
