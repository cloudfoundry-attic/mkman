package validators

import "fmt"

type emptinessValidator struct{}

func NewEmptinessValidator() Validator {
	return &emptinessValidator{}
}

func (e *emptinessValidator) Name() string {
	return "not be empty"
}

func (ev *emptinessValidator) Validate(vt ValidationTarget) error {
	if vt.object == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}
