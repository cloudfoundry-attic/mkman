package validators

// import (
// 	"fmt"
// 	"os"

// 	"github.com/cloudfoundry/multierror"
// )

// const (
// 	none     = 0
// 	FileType = 1 << iota
// 	DirType  = 1 << iota
// )

// type typeValidator struct {
// 	allowedType uint
// }

// func NewTypeValidator(allowedType uint) Validator {
// 	return &typeValidator{
// 		allowedType: allowedType,
// 	}
// }

// func (t *typeValidator) Validate(vt ValidationTarget) *multierror.MultiError {
// 	errors := multierror.NewMultiError(vt.name)
// 	fileInfo, _ := os.Stat(vt.object)
// 	if !t.isFileTypeAllowed(fileInfo) {
// 		errors.Add(fmt.Errorf("value must be absolute path to %s: '%s'", t.translate(), vt.object))
// 	}

// 	if errors.Length() > 0 {
// 		return errors
// 	}
// 	return nil
// }

// func (t *typeValidator) translate() string {
// 	switch t.allowedType {
// 	case FileType:
// 		return "file"
// 	case DirType:
// 		return "directory"
// 	case (FileType | DirType):
// 		return "file or directory"
// 	default:
// 		panic("unhandled")
// 	}
// }

// func (t typeValidator) isFileTypeAllowed(fileInfo os.FileInfo) bool {
// 	if fileInfo == nil {
// 		return false
// 	}

// 	if fileInfo.Mode().IsRegular() && (t.allowedType&FileType != none) {
// 		return true
// 	}

// 	if fileInfo.Mode().IsDir() && (t.allowedType&DirType != none) {
// 		return true
// 	}

// 	return false
// }
