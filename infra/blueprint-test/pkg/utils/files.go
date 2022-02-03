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
	_, err = f.Write([]byte(data))
	if err != nil {
		return "", err
	}
	f.Close()
	return f.Name(), nil
}
