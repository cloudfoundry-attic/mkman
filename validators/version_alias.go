package validators

import "fmt"

type versionAlias struct {
	versionAliases []string
}

func VersionAlias(versionAliases []string) Validator {
	return &versionAlias{
		versionAliases: versionAliases,
	}
}

func (v *versionAlias) ComposableName() string {
	return "valid version alias"
}

func (v *versionAlias) Validate(vt ValidationTarget) error {
	convertedObject, ok := vt.object.(string)
	if !ok {
		panic(fmt.Sprintf("Expected string type for %s", vt.name))
	}

	for _, element := range v.versionAliases {
		if element == convertedObject {
			return nil
		}
	}
	return fmt.Errorf("%s must be %s: %s", vt.name, v.ComposableName(), vt.object)
}
