package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/multierror"
)

type validator struct {
	name   string
	object string
	err    error
}

type Validation struct {
	VersionAliases *[]string
	AllowedType    uint
}

func NewValidator(object, name string) validator {
	return validator{
		name:   name,
		object: object,
	}
}

func (v *validator) Validate(validations Validation) error {
	err := v.validateNotEmpty()
	if err != nil {
		return err
	}
	if validations.VersionAliases != nil {
		err = v.validateVersionAlias(validations)
		if err == nil {
			return nil
		}
		err = v.validateIsAbsPath(validations)
		if err != nil {
			return err
		}
	}
	return v.validatePath(validations)
}

func (v *validator) validateNotEmpty() error {
	errors := multierror.NewMultiError(v.name)
	if v.object == "" {
		errors.Add(fmt.Errorf("value is required"))
		// Return when next error does not make sense
		return errors
	}
	return nil
}

func (v *validator) validatePath(validations Validation) error {

	errors := multierror.NewMultiError(v.name)

	err := v.validateIsAbsPath(validations)
	if err != nil {
		errors.Add(err)
		return errors
	}

	fileInfo, err := os.Stat(v.object)
	if os.IsNotExist(err) {
		errors.Add(fmt.Errorf("%s does not exist: '%s'", translate(validations.AllowedType), v.object))
		// Return when next error does not make sense
		return errors
	}

	if !isFileTypeAllowed(fileInfo, validations.AllowedType) {
		errors.Add(fmt.Errorf("value must be absolute path to %s: '%s'", translate(validations.AllowedType), v.object))
	}

	if errors.Length() > 0 {
		return errors
	}
	return nil
}

func (v *validator) validateIsAbsPath(validations Validation) error {
	if filepath.IsAbs(v.object) {
		return nil
	}

	if validations.VersionAliases != nil {
		return fmt.Errorf("value %s must be either a valid version alias or an absolute path", v.object)
	} else {
		return fmt.Errorf("value must be absolute path to %s: '%s'", translate(validations.AllowedType), v.object)
	}
}

func (v *validator) validateVersionAlias(validations Validation) error {
	for _, element := range *validations.VersionAliases {
		if element == v.object {
			return nil
		}
	}
	return fmt.Errorf("version alias not found")
}

func translate(allowedType uint) string {
	switch allowedType {
	case FileType:
		return "file"
	case DirType:
		return "directory"
	case (FileType | DirType):
		return "file or directory"
	default:
		panic("unhandled")
	}
}

func isFileTypeAllowed(fileInfo os.FileInfo, allowedPathType uint) bool {
	if fileInfo == nil {
		return false
	}

	if fileInfo.Mode().IsRegular() && (allowedPathType&FileType != none) {
		return true
	}

	if fileInfo.Mode().IsDir() && (allowedPathType&DirType != none) {
		return true
	}

	return false
}
