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
