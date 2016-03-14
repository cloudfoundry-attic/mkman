package validators

import (
	"fmt"
	"os"
	"strings"
)

type existsOnFilesystem struct {
	fileTypeValidators []FileTypeValidator
}

func ExistsOnFilesystem(fileTypeValidators ...FileTypeValidator) Validator {
	if len(fileTypeValidators) == 0 {
		panic(fmt.Errorf("must provide fileTypeValidator(s)"))
	}
	return &existsOnFilesystem{
		fileTypeValidators: fileTypeValidators,
	}
}

func (ev *existsOnFilesystem) ComposableName() string {
	allowableTypes := []string{}
	for _, v := range ev.fileTypeValidators {
		allowableTypes = append(allowableTypes, v.ComposableName())
	}
	errStr := strings.Join(allowableTypes, " or ")
	return fmt.Sprintf("a path to %s that exists", errStr)
}

func (ev *existsOnFilesystem) Validate(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	_, err := os.Stat(convertedObject)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s must be %s: '%s'", vt.name, ev.ComposableName(), vt.object)
	}
	var validators []Validator
	for _, fv := range ev.fileTypeValidators {
		validators = append(validators, fv.(Validator))
	}
	return AnyOf(validators...).Validate(vt)
}
