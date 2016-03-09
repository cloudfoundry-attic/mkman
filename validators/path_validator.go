package validators

import (
	"fmt"
	"os"

	//TODO Godep these paths
)

const (
	none     = 0
	FileType = 1 << iota
	DirType  = 1 << iota
)

type pathValidator struct {
	allowedType uint
}

func NewPathValidator(allowedType uint) Validator {
	return &pathValidator{
		allowedType: allowedType,
	}
}

func (pv *pathValidator) Name() string {
	return "valid path"
}

func (pv *pathValidator) Validate(vt ValidationTarget) error {
	fileInfo, err := os.Stat(vt.object)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: '%s'", pv.translate(), vt.object)
		// Return when next error does not make sense
	}

	if !pv.isFileTypeAllowed(fileInfo) {
		return fmt.Errorf("value must be path to %s: '%s'", pv.translate(), vt.object)
	}

	return nil
}

func (pv *pathValidator) translate() string {
	switch pv.allowedType {
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

func (pv pathValidator) isFileTypeAllowed(fileInfo os.FileInfo) bool {
	if fileInfo == nil {
		return false
	}

	if fileInfo.Mode().IsRegular() && (pv.allowedType&FileType != none) {
		return true
	}

	if fileInfo.Mode().IsDir() && (pv.allowedType&DirType != none) {
		return true
	}

	return false
}
