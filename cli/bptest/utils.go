package bptest

import (
	"io/fs"
	"io/ioutil"
	"path"
)

// findDirContentsFilter returns a list of files/directories in path based on a filter func
func findDirContentsFilter(dir string, filter func(fs.FileInfo) bool) ([]string, error) {
	results := make([]string, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return results, err
	}
	for _, f := range files {
		if filter(f) {
			results = append(results, path.Join(dir, f.Name()))
		}
	}
	return results, nil
}
