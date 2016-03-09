package validators

import (
	"fmt"
	"os"
)

type existenceValidator struct {
}

func NewExistenceValidator() Validator {
	return &existenceValidator{}
}

func (ev *existenceValidator) Name() string {
	return "exist"
}

func (ev *existenceValidator) Validate(vt ValidationTarget) error {
	_, err := os.Stat(vt.object)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: '%s'", vt.name, vt.object)
	}
	return nil
}
