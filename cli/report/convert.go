// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package report

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// ReadFilesAndConcat reads json files in a directory that are one object per row,
// and concats all objects into one single array
func ReadFilesAndConcat(dir string) (results []interface{}, err error) {
	files, err := listFiles(dir)
	const maxCapacity = 1024 * 1024
	if err != nil {
		return nil, err
	}

	for _, filePath := range files {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		buf := make([]byte, maxCapacity)
		s.Buffer(buf, maxCapacity)

		for s.Scan() {
			var row map[string]interface{}
			err = json.Unmarshal(s.Bytes(), &row)
			if err != nil {
				return nil, err
			}
			results = append(results, row)
		}
	}

	return
}

// listFiles returns a list of files under a dir. Errors will be grpc errors.
func listFiles(dir string) ([]string, error) {
	files := []string{}
	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "error visiting path %s", path)
		}
		if !f.IsDir() {
			files = append(files, path)
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		return nil, err
	}
	return files, nil
}
