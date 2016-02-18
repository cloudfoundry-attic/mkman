package tarball

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//go:generate counterfeiter . TarballReader
type TarballReader interface {
	ReadFile(filename string) ([]byte, error)
}

type tarballReader struct {
	tarPath string
}

func NewTarballReader(tarPath string) TarballReader {
	return &tarballReader{
		tarPath: tarPath,
	}
}

func (t *tarballReader) ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(t.tarPath)
	if err != nil {
		return nil, err
	}

	fileReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(fileReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			// We cannot test because getting tarReader to fail without EOF is hard
			return nil, err
		}

		if header.Name == filename {
			bytes, err := ioutil.ReadAll(tarReader)
			if err != nil {
				// We cannot test because getting ReadAll to fail is hard
				return nil, err
			}

			return bytes, nil
		}
	}

	return nil, fmt.Errorf(
		"filename '%s' not found in tarPath: %s\n",
		filename,
		t.tarPath,
	)
}
