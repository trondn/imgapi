package server

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func getSha1Sum(filename string) (sum string, err error) {
	// ok, generate the SHA1
	file, err := os.Open(filename)
	if err != nil {
		return sum, err
	}

	hasher := sha1.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return sum, err
	}

	sum = fmt.Sprintf("%x", hasher.Sum(nil))
	return sum, err
}
