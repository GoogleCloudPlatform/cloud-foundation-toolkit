package utils

import "os"

// WriteTmpFile writes data to a temp file and returns the path.
func WriteTmpFile(data string) (string, error) {
	f, err := os.CreateTemp("", "*")
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

// WriteTmpFileWithExtension writes data to a temp file with given extension and returns the path.
func WriteTmpFileWithExtension(data string, extension string) (string, error) {
	f, err := os.CreateTemp("", "*."+extension)
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
