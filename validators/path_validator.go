package validators

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/multierror" //TODO Godep these paths
)

type pathValidator struct {
	validations Validation
}

func NewPathValidator(validation Validation) Validator {
	return &pathValidator{
		validations: validation,
	}
}

func (pv *pathValidator) Validate(vt ValidationTarget) *multierror.MultiError {

	fmt.Println("In Path Validator ...")
	errors := multierror.NewMultiError(vt.name)

	err := pv.validateIsAbsPath(vt)
	if err != nil {
		errors.Add(err)
		return errors
	}

	fileInfo, err := os.Stat(vt.object)
	if os.IsNotExist(err) {
		errors.Add(fmt.Errorf("%s does not exist: '%s'", pv.translate(), vt.object))
		// Return when next error does not make sense
		return errors
	}

	if !pv.isFileTypeAllowed(fileInfo) {
		errors.Add(fmt.Errorf("value must be absolute path to %s: '%s'", pv.translate(), vt.object))
	}

	if errors.Length() > 0 {
		return errors
	}
	return nil
}

func (pv *pathValidator) translate() string {
	switch pv.validations.AllowedType {
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

func (pv pathValidator) validateIsAbsPath(vt ValidationTarget) error {
	if filepath.IsAbs(vt.object) {
		return nil
	}

	if pv.validations.VersionAliases != nil {
		return fmt.Errorf("value %s must be either a valid version alias or an absolute path", vt.object)
	} else {
		return fmt.Errorf("value must be absolute path to %s: '%s'", pv.translate(), vt.object)
	}
}

func (pv pathValidator) isFileTypeAllowed(fileInfo os.FileInfo) bool {
	if fileInfo == nil {
		return false
	}

	if fileInfo.Mode().IsRegular() && (pv.validations.AllowedType&FileType != none) {
		return true
	}

	if fileInfo.Mode().IsDir() && (pv.validations.AllowedType&DirType != none) {
		return true
	}

	return false
}
