package validators

import "fmt"

type nonEmpty struct{}

func NonEmpty() Validator {
	return &nonEmpty{}
}

func (e *nonEmpty) ComposableName() string {
	return "non-empty"
}

func (ev *nonEmpty) Validate(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	if convertedObject == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}
