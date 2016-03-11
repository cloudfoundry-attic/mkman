package jobtemplates

import "github.com/cloudfoundry/mkman/stubmakers"

type jobTemplateStub struct {
	Meta meta `yaml:"meta"`
}

type meta struct {
	EtcdRelease etcdRelease `yaml:"etcd_release,omitempty"`
	Merge       string      `yaml:"<<"`
}

type etcdRelease struct {
	Name string `yaml:"name"`
}

type etcdTemplate struct {
	Name    string `yaml:"name"`
	Release string `yaml:"release"`
}

type jobTemplateStubMaker struct{}

func NewJobTemplateStubMaker() stubmakers.StubMaker {
	return &jobTemplateStubMaker{}
}

func (t *jobTemplateStubMaker) MakeStub() (string, error) {
	jobTemplate := jobTemplateStub{
		Meta: meta{
			Merge: "(( merge ))",
			EtcdRelease: etcdRelease{
				Name: "etcd",
			},
		},
	}
	return stubmakers.MarshalTempStub(&jobTemplate, "job-template.yml")
}
