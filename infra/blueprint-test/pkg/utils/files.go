package utils

import (
	"io/ioutil"
)

// WriteTmpFile writes data to a temp file and returns the path.
func WriteTmpFile(data string) (string, error) {
	f, err := ioutil.TempFile("", "*")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
