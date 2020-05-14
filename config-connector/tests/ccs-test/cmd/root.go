// Copyright 2020 Google LLC
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

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/config-connector/tests/ccs-test/util"
	"github.com/spf13/cobra"
)

const (
	kubectlBinaryName          = "kubectl"
	testsDirPath               = "config-connector/tests"
	originalValuesFileName     = "original_values.yaml"
	requiredFieldsOnlyFileName = "required_fields_only.yaml"
	envFileRelativePath        = "testcases/environments.yaml"
	yamlFileSuffix             = ".yaml"
)

var (
	relativePath string
	timeout string

	// Regex of the env vars that require a randomized suffix.
	// It should be in the format of $ENV_VAR-$RANDOM_ID.
	re = regexp.MustCompile(`\$(?P<EnvName>[A-Z]+|[A-Z]+[A-Z_]*[A-Z]+)(-\$RANDOM_ID)`)

	rootCmd = &cobra.Command{
		Use:   "ccs-test",
		Short: "CLI to test Config Connector Solutions",
		Long:  `CLI to test Config Connector Solutions`,
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a given test by the relative path --path",
		Long:  "Run a given test by the relative path --path",
		Run: func(cmd *cobra.Command, args []string) {
			// Check the required flag.
			if relativePath == "" {
				log.Fatal("\"--path\" must be specified to run test")
			}

			log.Printf("======Testing solution %q...======\n", relativePath)

			// Calculate the path to testcase directory and solution directory.
			current, err := os.Getwd()
			if err != nil {
				log.Fatalf("error retrieving the current directory: %v", err)
			}
			if !strings.HasSuffix(current, testsDirPath) {
				log.Fatalf("error running tests under directory: %s. Please "+
					"follow the instructions in the README.", current)
			}
			testCasePath := filepath.Join(current, "testcases", relativePath)
			envFilePath := filepath.Join(current, envFileRelativePath)

			parent := filepath.Dir(current)
			solutionPath := filepath.Join(parent, "solutions", relativePath)

			// Clean up the left over resources if there are any.
			if err := deleteResources(solutionPath); err != nil {
				log.Fatalf("error cleaning up resources before running the "+
					"test. Please clean them up manually: %v", err)
			}

			// Fetch the testcase values and run the test.
			envValues := make(map[string]string)
			if err := parseYamlToStringMap(envFilePath, envValues); err != nil {
				log.Fatalf("error retrieving envrionment variables: %v", err)
			}

			originalValues := make(map[string]string)
			if err := parseYamlToStringMap(filepath.Join(testCasePath, originalValuesFileName), originalValues); err != nil {
				log.Fatalf("error retrieving orginal values: %v", err)
			}

			testValues := make(map[string]string)
			if err := parseYamlToStringMap(filepath.Join(testCasePath, requiredFieldsOnlyFileName), testValues); err != nil {
				log.Fatalf("error retrieving test values: %v", err)
			}

			// Generate the random IDs first.
			randomId, err := util.GenerateRandomizedSuffix()
			if err != nil {
				log.Fatalf("error generating the randomized suffix for resource names: %v", err)
			}

			// Then populate the values of the env var and the random id.
			realValues, err := finalizeValues(randomId, envValues, testValues)
			if err != nil {
				log.Fatalf("error finalizing test values: %v", err)
			}

			if err := runKptTestcase(solutionPath, timeout, realValues, originalValues); err != nil {
				log.Fatalf("test failed for solution %q: %v", relativePath, err)
			}

			log.Printf("======Successfully finished the test for solution %q======\n", relativePath)
		},
	}
)

func init() {
	runCmd.PersistentFlags().StringVarP(&relativePath, "path", "p", "", "[Required] The relative path to the folder of the solution's testcases, e.g. `iam/kpt/member-iam`.")
	runCmd.PersistentFlags().StringVarP(&timeout, "timeout", "t", "60s", "[Optional] The timeout used to wait for resources to be READY. Default: `60s`.")
	rootCmd.AddCommand(runCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseYamlToStringMap(filePath string, result map[string]string) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file '%s': %v", filePath, err)
	}
	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		return fmt.Errorf("error unmarshaling file '%s': %v", filePath, err)
	}
	return nil
}

func finalizeValues(randomId string, envValues map[string]string, testValues map[string]string) (map[string]string, error) {
	realValues := make(map[string]string)
	for key, value := range testValues {
		if !strings.HasPrefix(value, "$") {
			return nil, fmt.Errorf("test value for setter %q is %q, expect a reference, e.g. $ENV_VAR", key, value)
		}
		realValue := ""
		ok := false
		if re.MatchString(value) {
			submatch := re.FindStringSubmatch(value)
			if len(submatch) == 0 {
				return nil, fmt.Errorf("env var name is invalid in test value %q", value)
			}
			subexpNames := re.SubexpNames()
			for i, name := range subexpNames {
				if name == "EnvName" {
					prefix, ok := envValues[submatch[i]]
					if !ok {
						return nil, fmt.Errorf("couldn't find the env var %q", submatch[i])
					}
					realValue = fmt.Sprintf("%s-%s", prefix, randomId)
					break
				}
			}
		} else {
			realValue, ok = envValues[strings.TrimPrefix(value, "$")]
			if !ok {
				return nil, fmt.Errorf("couldn't find the env var %q", strings.TrimPrefix(value, "$"))
			}
		}
		realValues[key] = realValue
	}

	return realValues, nil
}

func runKptTestcase(solutionPath string, timeout string, testValues map[string]string, originalValues map[string]string) error {
	// Set the kpt setters defined in the testcase.
	log.Println("======Setting the kpt setters...======")
	for key, value := range testValues {
		output, err := exec.Command("kpt", "cfg", "set", solutionPath, key,
			value, "--set-by", "test").CombinedOutput()
		if err != nil {
			log.Printf("stderr:\n%v\nstdout:\n%s\n", err, string(output))
			errToReturn := fmt.Errorf("error setting setter '%s' with value "+
				"'%s': %v\nstdout: %s", key, value, err, string(output))

			// Clean up before exit with errors.
			if err := resetKptSetters(solutionPath, originalValues); err != nil {
				return concatErrors(
					"error resetting kpt setters before exiting",
					err, errToReturn)
			}

			return errToReturn
		}
		log.Printf("%s\n", string(output))
	}
	log.Println("======Successfully set the kpt setters======")

	// Apply all the resources.
	log.Println("======Creating the resources...======")
	output, err := exec.Command("kubectl", "create", "-f", solutionPath).CombinedOutput()
	if err != nil {
		log.Printf("stderr:\n%v\nstdout:\n%s\n", err, string(output))
		errToReturn := fmt.Errorf("error creating resources: %v\nstdout: %s", err, string(output))

		// Clean up before exit with errors.
		if err := cleanUp(solutionPath, originalValues); err != nil {
			return concatErrors(
				"error cleanning up resources before exiting",
				err, errToReturn)
		}
		return errToReturn
	}
	log.Printf("%s\n", string(output))
	log.Println("======Successfully created the resources======")

	// Wait for all the resources to be ready.
	if err := verifyReadyCondition(solutionPath, timeout); err != nil {
		errToReturn := fmt.Errorf("error verifying the ready condition: %v", err)

		// Clean up before exit with errors.
		if err := cleanUp(solutionPath, originalValues); err != nil {
			return concatErrors(
				"error cleanning up resources before exiting",
				err, errToReturn)
		}
		return errToReturn
	}

	// Clean up.
	return cleanUp(solutionPath, originalValues)
}

func cleanUp(solutionPath string, originalValues map[string]string) error {
	resourceErr := deleteResources(solutionPath)
	setterErr := resetKptSetters(solutionPath, originalValues)
	if resourceErr != nil || setterErr != nil {
		return concatErrors(
			fmt.Sprintf("error cleanning up the test for solution %q. "+
				"Please manually delete the Config Connector resources and "+
				"reset the kpt setters", solutionPath),
			resourceErr, setterErr)
	}
	return nil
}

func verifyReadyCondition(solutionPath string, timeout string) error {
	log.Println("======Verifying that all the Config Connector resources are ready...======")

	files, err := ioutil.ReadDir(solutionPath)
	if err != nil {
		return fmt.Errorf("error reading solution directory %q: %v", solutionPath, err)
	}

	for _, file := range files {
		// We should only verify the YAML config files for Config Connector
		// resources.
		fileName := file.Name()
		if !strings.HasSuffix(fileName, yamlFileSuffix) || strings.Contains(fileName, "namespace") {
			continue
		}

		resourceFilePath := filepath.Join(solutionPath, fileName)
		output, err := exec.Command("kubectl", "wait", "--for=condition=ready",
			"-f", resourceFilePath, fmt.Sprintf("--timeout=%s", timeout)).CombinedOutput()
		if err != nil {
			log.Printf("stderr:\n%v\nstdout:\n%s\n", err, string(output))
			errToReturn := fmt.Errorf("resource in file %q is not ready in timeout: %v\nstdout: %s", fileName, timeout, err, string(output))
			status, err := getSolutionResourceStatus(solutionPath)
			if err != nil {
				return concatErrors("error printing resource status", err, errToReturn)
			}

			return fmt.Errorf("%v\nResource status:\n%s", errToReturn, status)
		}
		log.Printf("%s\n", string(output))

	}

	log.Println("======All the Config Connector resrouces are ready======")
	return nil
}

func getSolutionResourceStatus(solutionPath string) (string, error) {
	output, err := exec.Command("kubectl", "get", "-f", solutionPath,
		"-o=custom-columns=NAME:.metadata.name,KIND:.kind,CONDITION.REASON:.status.conditions[0].reason,CONDITION.MESSAGE:.status.conditions[0].message").
		CombinedOutput()

	if err != nil {
		log.Printf("stderr:\n%v\nstdout:\n%s\n", err, string(output))
		return "", fmt.Errorf("error getting the status of the resource(s): %v\nstdout: %s", err, string(output))
	}

	return string(output), nil
}

func deleteResources(solutionPath string) error {
	log.Println("======Deleting the resources...======")
	output, err := exec.Command("kubectl", "delete", "-f", solutionPath, "--wait").CombinedOutput()
	if err != nil {
		log.Printf("stderr:\n%v\nstdout:\n%s\n", err, string(output))
		err = fmt.Errorf("error deleting resources: %v\nstdout: %s", err, string(output))
		if isNotFoundErrorOnly(err) {
			log.Println(err)
			log.Println("======Finished deleting the resources======")
			return nil
		}
		return fmt.Errorf("error deleting resources: %v\nstdout: %s", err, string(output))
	}
	log.Printf("%s\n", output)
	log.Println("======Successfully deleted the resources======")
	return nil
}

func isNotFoundErrorOnly(err error) bool {
	numErrors := strings.Count(err.Error(), "Error from server")
	numNotFoundErrors := strings.Count(err.Error(), "Error from server (NotFound)")
	return numErrors == numNotFoundErrors
}

func resetKptSetters(solutionPath string, originalValues map[string]string) error {
	log.Println("======Resetting the kpt setters...======")
	for key, value := range originalValues {
		output, err := exec.Command("kpt", "cfg", "set", solutionPath, key, value, "--set-by", "PLACEHOLDER").CombinedOutput()
		if err != nil {
			log.Printf("stderr:\n%v\nstdout:\n%s\n", err, string(output))
			return fmt.Errorf("error setting setter '%s' back to the original value '%s': %v\nstdout: %s", key, value, err, string(output))
		}
		log.Printf("%s\n", string(output))
	}
	log.Println("======Successfully reset the kpt setters======")
	return nil
}

func concatErrors(msg string, errs ...error) error {
	errToReturn := errors.New(msg)
	for _, err := range errs {
		if err == nil {
			continue
		}
		errToReturn = fmt.Errorf("%v:%v", errToReturn, err)
	}
	return errToReturn
}
