package validators

type allOf struct {
	validators  []Validator
	firstFailed Validator
}

func AllOf(validators ...Validator) Validator {
	return &allOf{validators: validators}
}

func (a *allOf) ComposableName() string {
	if a.firstFailed == nil {
		return ""
	}
	return a.firstFailed.ComposableName()
}

func (a *allOf) Validate(vt ValidationTarget) error {
	for _, validator := range a.validators {
		err := validator.Validate(vt)
		if err != nil {
			a.firstFailed = validator
			return err
		}
	}
	return nil
}
