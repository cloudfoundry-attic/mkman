package validators

import "fmt"

type versionAliasValidator struct {
	versionAliases []string
}

func VersionAliasValidator(versionAliases []string) Validator {
	return &versionAliasValidator{
		versionAliases: versionAliases,
	}
}

func (v *versionAliasValidator) ComposableName() string {
	return "valid version alias"
}

func (v *versionAliasValidator) Validate(vt ValidationTarget) error {
	for _, element := range v.versionAliases {
		convertedObject, ok := vt.object.(string)
		if !ok {
			panic(fmt.Sprintf("Expected string type for %s", vt.name))
		}

		if element == convertedObject {
			return nil
		}
	}
	return fmt.Errorf("version alias not found")
}
