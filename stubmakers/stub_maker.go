package stubmakers

//go:generate counterfeiter . StubMaker
type StubMaker interface {
	MakeStub() (string, error)
}
