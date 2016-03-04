package testhelpers

import "fmt"

type AlwaysErrorWriter struct{}

func (w *AlwaysErrorWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("writer error")
}
