package validators

import "fmt"

type andValidator struct {
	validators  []Validator
	firstFailed Validator
}

func And(validators ...Validator) Validator {
	return &andValidator{validators: validators}
}

func (a *andValidator) Name() string {
	return a.firstFailed.Name()
}

func (a *andValidator) Validate(vt ValidationTarget) error {
	for _, v := range a.validators {
		err := v.Validate(vt)
		if err != nil {
			a.firstFailed = v
			return fmt.Errorf("%s", v.Name())
		}
	}
	return nil
}
