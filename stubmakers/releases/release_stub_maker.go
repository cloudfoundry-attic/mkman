package releases

import (
	"github.com/cloudfoundry/mkman/stubmakers"
	"github.com/cloudfoundry/mkman/stubmakers/releases/releasemakers"
)

type ReleaseStub struct {
	Releases []releasemakers.Release `yaml:"releases,omitempty"`
}

type releaseStubMaker struct {
	releaseMakers []releasemakers.ReleaseMaker
}

func NewReleaseStubMaker(r []releasemakers.ReleaseMaker) stubmakers.StubMaker {
	return &releaseStubMaker{
		releaseMakers: r,
	}
}

func (r *releaseStubMaker) MakeStub() (string, error) {
	stub := ReleaseStub{}
	for _, releaseMaker := range r.releaseMakers {
		release, err := releaseMaker.MakeRelease()
		if err != nil {
			return "", err
		}
		stub.Releases = append(stub.Releases, *release)
	}

	return stubmakers.MarshalTempStub(stub, "release.yml")
}
