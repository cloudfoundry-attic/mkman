package jobtemplates

import "github.com/cloudfoundry/mkman/stubmakers"

type jobTemplateStub struct {
	Meta meta `yaml:"meta"`
}

type meta struct {
	EtcdTemplates []etcdTemplate `yaml:"etcd_templates,omitempty"`
	Merge         string         `yaml:"<<"`
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
			EtcdTemplates: []etcdTemplate{
				etcdTemplate{
					Name:    "etcd",
					Release: "etcd",
				},
				etcdTemplate{
					Name:    "etcd_metrics_server",
					Release: "etcd",
				},
			},
		},
	}
	return stubmakers.MarshalTempStub(&jobTemplate, "job-template.yml")
}
