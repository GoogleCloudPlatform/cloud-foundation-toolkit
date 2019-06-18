package deployment

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// The runGCloud function runs the gcloud tool with the specified arguments. It is implemented
// as a variable so that it can be mocked in tests of its exported consumers.
var runGCloud = func(args ...string) (result string, err error) {
	args = append(args, "--format", "yaml")
	log.Println("gcloud", strings.Join(args, " "))
	cmd := exec.Command("gcloud", args...)
	outBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}
	cmd.Stdout = outBuff
	cmd.Stderr = errBuff
	// pass user's PATH env variable, expected gcloud executable can be found in PATH
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start cmd: %v, \n Output:\n %v", err, errBuff.String())
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Cmd returned error: %v, \n Output:\n %v", err, errBuff.String())
		return errBuff.String(), err
	}

	return outBuff.String(), err
}

// GetOutputs retrive existing Deployment outputs using gcloud and store result in map[string]string
// where "resourceName.propertyName" is key, and value is string representation of the output value.
func GetOutputs(name string, project string) (map[string]string, error) {
	data, err := runGCloud("deployment-manager", "manifests", "describe", "--deployment", name, "--project", project)
	if err != nil {
		log.Printf("Failed to get deployment manifest: %v", err)
		return nil, err
	}
	return parseOutputs(data)
}

// Create creates a new Deployment based on a Deployment object passed into it.
func Create(deployment *Deployment) (*Deployment, error) {
	args := []string{
		"deployment-manager",
		"deployments",
		"create",
		deployment.config.Name,
		"--config",
		deployment.configFile,
		"--project",
		deployment.config.Project,
	}
	_, err := runGCloud(args...)
	if err != nil {
		log.Printf("Failed to create deployment: %v, error: %v", deployment, err)
		return nil, err
	}
	outputs, err := GetOutputs(deployment.config.Name, deployment.config.Project)
	if err != nil {
		log.Printf("Failed to get outputs for deployment: %v, error: %v", deployment, err)
		return nil, err
	}
	deployment.Outputs = outputs
	return deployment, nil
}

func parseOutputs(data string) (map[string]string, error) {
	describe, err := unmarshal(data)
	if err != nil {
		log.Println("Error parsing deployment outputs")
		return nil, err
	}

	layoutData := describe["layout"].(string)

	res := &struct {
		Resources []struct {
			Name    string
			Outputs []struct {
				Value interface{} `yaml:"finalValue"`
				Name  string
			}
		}
	}{}
	err = yaml.Unmarshal([]byte(layoutData), res)
	if err != nil {
		log.Println("Error parsing deployment outputs layout section")
		return nil, err
	}

	result := make(map[string]string)
	for _, resource := range res.Resources {
		for _, output := range resource.Outputs {
			key := resource.Name + "." + output.Name
			switch value := output.Value.(type) {
			case string:
				result[key] = value
			case map[interface{}]interface{}:
				log.Println(key + " is map")
			}
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}
