package validators

import "fmt"

type versionAliasValidator struct {
	versionAliases []string
}

func NewVersionAliasValidator(versionAliases []string) Validator {
	return &versionAliasValidator{
		versionAliases: versionAliases,
	}
}

func (v *versionAliasValidator) Name() string {
	return "valid version alias"
}

func (v *versionAliasValidator) Validate(vt ValidationTarget) error {
	for _, element := range v.versionAliases {
		if element == vt.object {
			return nil
		}
	}
	return fmt.Errorf("version alias not found")
}
