package validators

import "fmt"

type nonEmptinessValidator struct{}

func NonEmptinessValidator() Validator {
	return &nonEmptinessValidator{}
}

func (e *nonEmptinessValidator) ComposableName() string {
	return "non-empty"
}

func (ev *nonEmptinessValidator) Validate(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	if convertedObject == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}
