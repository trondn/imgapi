package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

/**
 * Utility function to get the SHA1 sum for a named file
 *
 * @param filename the name of the file to read
 * @return sum The SHA1 sum of the file in ASCII
 *         err The error object if something failed
 */
func GetSha1Sum(filename string) (sum string, err error) {
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
