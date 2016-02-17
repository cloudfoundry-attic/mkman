package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func ReadFileContentsFromTar(tarPath string, filename string) ([]byte, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		panic(err)
	}

	fileReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}

	tr := tar.NewReader(fileReader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if hdr.Name == filename {
			b, err := ioutil.ReadAll(tr)
			if err != nil {
				panic(err)
			}

			return b, nil
		}
	}

	return nil, fmt.Errorf(
		"filename: %s not found in tarPath: %s\n",
		filename,
		tarPath,
	)
}
