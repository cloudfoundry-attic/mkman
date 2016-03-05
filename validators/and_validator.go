package validators

import "fmt"

type andValidator struct {
	validators  []Validator
	firstFailed Validator
}

func And(validators ...Validator) Validator {
	return &andValidator{validators: validators}
}

func (a *andValidator) ComposableName() string {
	return a.firstFailed.ComposableName()
}

func (a *andValidator) Validate(vt ValidationTarget) error {
	for _, v := range a.validators {
		err := v.Validate(vt)
		if err != nil {
			a.firstFailed = v
			return fmt.Errorf("value must be %s: %s", v.ComposableName(), vt.object)
		}
	}
	return nil
}
