package validators

import "fmt"

type nonEmptyArray struct{}

func NonEmptyArray() Validator {
	return &nonEmptyArray{}
}

func (e *nonEmptyArray) ComposableName() string {
	return "non-empty array"
}

func (ev *nonEmptyArray) Validate(vt ValidationTarget) error {
	switch v := vt.object.(type) {
	case []string:
		if len(v) <= 0 {
			return fmt.Errorf("value must be %s: %s", ev.ComposableName(), vt.name)
		}
		return nil
	default:
		return fmt.Errorf("value must be of type string array: %s", vt.name)
	}
}
