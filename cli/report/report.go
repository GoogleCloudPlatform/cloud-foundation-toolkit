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
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
)

// GenerateReports takes raw CAI exports from <dirPath> directory,
// run rego queries defined <queryPath> directory,
// and generate output of <reportFormat> in <outputPath> directory
func GenerateReports(dirPath string, queryPath string, outputPath string, reportFormat string) error {
	fileSuffix := time.Now().Format("2006.01.02-15.04.05")
	rawAssetFileName, err := convertAndGenerateTempAssetFile(dirPath, outputPath, fileSuffix)
	if err != nil {
		return err
	}
	results, err := generateReportData(rawAssetFileName, queryPath, outputPath)
	if err != nil {
		return err
	}
	printReports(results, outputPath, reportFormat, fileSuffix)
	return nil
}

func convertAndGenerateTempAssetFile(caiPath string, outputPath string, fileMidName string) (rawAssetFileName string, err error) {
	results, err := ReadFilesAndConcat(caiPath)
	if err != nil {
		return "", err
	}
	wrapped := map[string]interface{}{
		"assets": results,
	}
	outJSON, _ := json.MarshalIndent(wrapped, "", "  ")
	rawAssetFileName = "raw_assets_" + fileMidName + ".json"
	err = ioutil.WriteFile(filepath.Join(outputPath, rawAssetFileName), outJSON, 0644)
	if err != nil {
		return "", err
	}
	return
}

func findReports(paths []string) (results interface{}, err error) {
	// Load resources from json and rego files
	resources, err := loader.All(paths)
	if err != nil {
		return nil, err
	}
	compiler, err := resources.Compiler()
	if err != nil {
		return nil, err
	}
	store, err := resources.Store()
	if err != nil {
		return nil, err
	}
	r := rego.New(
		rego.Query(`data.reports`),
		rego.Compiler(compiler),
		rego.Store(store),
	)
	rs, err := r.Eval(context.Background())
	if err != nil {
		return nil, err
	}
	results = rs[0].Expressions[0].Value
	return results, err
}

func generateReportData(rawAssetFileName string, queryPath string, outputPath string) (results interface{}, err error) {
	return findReports([]string{filepath.Join(outputPath, rawAssetFileName), queryPath})
}

func printReports(results interface{}, reportOutputPath string, format string, fileSuffix string) error {
	resultsMap := results.(map[string]interface{})
	for group, reports := range resultsMap {
		reportsMap := reports.(map[string]interface{})
		for reportName, content := range reportsMap {
			if strings.HasSuffix(reportName, "_report") {
				reportFileName := group + "." + reportName + "_" + fileSuffix
				fmt.Printf("Generating %v.%v\n", group, reportName)
				if format == "json" {
					reportFileName = reportFileName + ".json"
					fileContent, err := json.MarshalIndent(content, "", "  ")
					if err != nil {
						return err
					}
					err = ioutil.WriteFile(filepath.Join(reportOutputPath, reportFileName), fileContent, 0644)
					if err != nil {
						return err
					}
				} else {
					reportFileName = reportFileName + ".csv"
					contentSlice := content.([]interface{})
					f, _ := os.Create(filepath.Join(reportOutputPath, reportFileName))

					defer f.Close()
					w := csv.NewWriter(f)
					if len(contentSlice) > 0 {
						firstRow := contentSlice[0]
						var keys []string
						firstRowMap := firstRow.(map[string]interface{})
						for key := range firstRowMap {
							keys = append(keys, key)
						}
						sort.Strings(keys)
						w.Write(keys)
						w.Flush()
						for _, record := range contentSlice {
							recordMap := record.(map[string]interface{})
							var record []string
							for _, key := range keys {
								record = append(record, recordMap[key].(string))
							}
							w.Write(record)
						}
						w.Flush()
					}
				}
			}
		}
	}
	return nil
}

// ListAvailableReports lists the names of available reports in queryPath
func ListAvailableReports(queryPath string) error {
	results, error := findReports([]string{queryPath})

	resultsMap := results.(map[string]interface{})
	for group, reports := range resultsMap {
		reportsMap := reports.(map[string]interface{})
		for reportName := range reportsMap {
			if strings.HasSuffix(reportName, "_report") {
				fmt.Println(group + "." + reportName)
			}
		}
	}

	return error
}
