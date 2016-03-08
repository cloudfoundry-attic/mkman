package validators

import "github.com/cloudfoundry/multierror"

// create a validator on an object
//    target = NewValidationTarget(FPATH, "cf")
//		target.AddValidators([]Validator {
//			NotEmptyValidator(),
//			IsAbsPathValidator(),
//			IsAliasValidator(AllowedAliasList),
//			IsOfCorrectTypeValidator(AllowedTypeList)
//		})
//		errs := target.Validate()
//
//		func (v ValidationTarget) Validate() error {
//			for _, validation := range v.Validators {
//				validation.validate(v)
//			}
//		}
//

const (
	none     = 0
	FileType = 1 << iota
	DirType  = 1 << iota
)

type ValidationTarget struct {
	name       string
	object     string
	validators *[]Validator
}

func NewValidationTarget(object, name string) ValidationTarget {
	return ValidationTarget{
		name:       name,
		object:     object,
		validators: &[]Validator{},
	}
}

func (vt ValidationTarget) Validate() *multierror.MultiError {
	errs := multierror.NewMultiError(vt.name)
	for _, v := range *vt.validators {
		errs.Add(v.Validate(vt))
	}
	return errs
}

func (vt ValidationTarget) Add(validators []Validator) {
	for _, v := range validators {
		*vt.validators = append(*vt.validators, v)
	}
}

type Validator interface {
	Validate(vt ValidationTarget) *multierror.MultiError
}

type Validation struct {
	VersionAliases *[]string
	AllowedType    uint
}

// func (v *validator) Validate(validations Validation) error {
// 	err := v.validateNotEmpty()
// 	if err != nil {
// 		return err
// 	}
// 	if validations.VersionAliases != nil {
// 		err = v.validateVersionAlias(validations)
// 		if err == nil {
// 			return nil
// 		}
// 		err = v.validateIsAbsPath(validations)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return v.validatePath(validations)
// }

// func (v *validator) validateVersionAlias(validations Validation) error {
// 	for _, element := range *validations.VersionAliases {
// 		if element == v.object {
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("version alias not found")
// }
