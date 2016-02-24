package multierror

import "fmt"

type MultiError struct {
	errors []error
}

func (e MultiError) Error() string {
	var errStr string
	if len(e.errors) == 1 {
		errStr = "encountered 1 error during validation:\n"
	} else {
		errStr = fmt.Sprintf("encountered %d errors during validation:\n", len(e.errors))
	}

	for _, err := range e.errors {
		errStr = fmt.Sprintf("%s%s\n", errStr, err.Error())
	}
	return errStr
}

// Add an error to the collection of errors.
// err must be non-nil
func (e *MultiError) Add(err error) {
	errors, ok := err.(MultiError)
	if ok {
		e.errors = append(e.errors, errors.errors...)
	} else {
		e.errors = append(e.errors, err)
	}
}

// Add an error to the collection of errors with a provided prefix.
// err must be non-nil
// prefix can be empty
func (e *MultiError) AddWithPrefix(err error, prefix string) {
	errors, ok := err.(MultiError)
	if ok {
		errors.prefixAll(prefix)
		e.errors = append(e.errors, errors.errors...)
	} else {
		e.errors = append(e.errors, fmt.Errorf("%s%s", prefix, err.Error()))
	}
}

func (e *MultiError) prefixAll(prefix string) {
	for i, err := range e.errors {
		e.errors[i] = fmt.Errorf("%s%s", prefix, err.Error())
	}
}

func (e *MultiError) HasAny() bool {
	return len(e.errors) > 0
}
