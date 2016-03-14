package validators

import "fmt"

type nonEmptyArrayValidator struct{}

func NonEmptyArrayValidator() Validator {
	return &nonEmptyArrayValidator{}
}

func (e *nonEmptyArrayValidator) ComposableName() string {
	return "non-empty"
}

func (ev *nonEmptyArrayValidator) Validate(vt ValidationTarget) error {
	switch vt.object.(type) {
	case []string:
		v := vt.object.([]string)
		if len(v) <= 0 {
			return fmt.Errorf("value must be non-empty array: %s", vt.name)
		}
	default:
		return fmt.Errorf("value must of type string array: %s", vt.name)
	}

	return nil
}
