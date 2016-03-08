package validators

import (
	"fmt"

	"github.com/cloudfoundry/multierror"
)

type emptinessValidator struct{}

func NewEmptinessValidator() Validator {
	return &emptinessValidator{}
}

func (ev *emptinessValidator) Validate(vt ValidationTarget) *multierror.MultiError {
	errors := multierror.NewMultiError(vt.name)
	if vt.object == "" {
		errors.Add(fmt.Errorf("value is required"))
		// Return when next error does not make sense
		return errors
	}
	return nil
}
