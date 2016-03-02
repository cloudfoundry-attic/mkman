package releasemakers

type cfReleaseMaker struct {
	releasePath string
}

func NewCfReleaseMaker(releasePath string) ReleaseMaker {
	return &cfReleaseMaker{
		releasePath: releasePath,
	}
}

func (r *cfReleaseMaker) MakeRelease() (*Release, error) {
	return &Release{
		Name:    "cf",
		URL:     "file://" + r.releasePath,
		Version: "create",
	}, nil
}
