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
	fmt.Printf("@@@ DEBUG opening tar: %s\n", tarPath)
	file, err := os.Open(tarPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("@@@ DEBUG unzipping: %s\n", tarPath)
	fileReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("@@@ DEBUG untarring: %s\n", tarPath)
	tr := tar.NewReader(fileReader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Printf("@@@ DEBUG name: %s\n", hdr.Name)
		if hdr.Name == filename {
			fmt.Printf("@@@ DEBUG found: %s\n", filename)
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
