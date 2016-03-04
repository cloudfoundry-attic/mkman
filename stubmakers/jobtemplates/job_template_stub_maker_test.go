package jobtemplates_test

import (
	"io/ioutil"
	"regexp"
	"strings"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/cloudfoundry/mkman/stubmakers"
	"github.com/cloudfoundry/mkman/stubmakers/jobtemplates"
)

var _ = Describe("TemplateStubMaker", func() {
	var templateStubMaker stubmakers.StubMaker

	BeforeEach(func() {
		templateStubMaker = jobtemplates.NewJobTemplateStubMaker()
	})

	It("writes a release stub and returns the path", func() {
		stubPath, err := templateStubMaker.MakeStub()
		Expect(err).NotTo(HaveOccurred())

		fileBytes, err := ioutil.ReadFile(stubPath)
		Expect(err).NotTo(HaveOccurred())

		remainder, mergeLines := splitOutMergeLine(fileBytes)
		Expect(mergeLines).To(HaveLen(1))
		Expect(mergeLines[0]).To(ContainSubstring("<<: (( merge ))"))

		var templateStub struct {
			Meta struct {
				EtcdTemplates []struct {
					Name    string `yaml:"name"`
					Release string `yaml:"release"`
				} `yaml:"etcd_templates"`
			} `yaml:"meta"`
		}

		remainderBytes := []byte(strings.Join(remainder, "\n"))

		err = yaml.Unmarshal(remainderBytes, &templateStub)
		Expect(err).NotTo(HaveOccurred())

		Expect(templateStub.Meta.EtcdTemplates).To(HaveLen(2))

		Expect(templateStub.Meta.EtcdTemplates[0].Name).To(Equal("etcd"))
		Expect(templateStub.Meta.EtcdTemplates[0].Release).To(Equal("etcd"))
		Expect(templateStub.Meta.EtcdTemplates[1].Name).To(Equal("etcd_metrics_server"))
		Expect(templateStub.Meta.EtcdTemplates[1].Release).To(Equal("etcd"))
	})
})

func splitOutMergeLine(fileBytes []byte) (nonMerge, merge []string) {
	mergeRegex := regexp.MustCompile("merge")

	fileLines := strings.Split(string(fileBytes), "\n")
	for _, line := range fileLines {
		if mergeRegex.Match([]byte(line)) {
			merge = append(merge, line)
		} else {
			nonMerge = append(nonMerge, line)
		}
	}
	return
}
