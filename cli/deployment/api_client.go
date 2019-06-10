package deployment

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type Status int

const (
	Done     Status = 0
	Pending  Status = 1
	Running  Status = 2
	NotFound Status = 3
	Error    Status = -1
)

const (
	ActionApply  string = "apply"
	ActionDelete string = "delete"
	ActionCreate string = "create"
	ActionUpdate string = "update"
)

var actions = []string{ActionApply, ActionDelete, ActionCreate, ActionUpdate}

// Function runGCloud exposed to variable in order to mock it inside api_client_test.go
// The runGCloud function runs the gcloud tool with the specified arguments. It is implemented
// as a variable so that it can be mocked in tests of its exported consumers.
var runGCloud = func(args ...string) (result string, err error) {
	args = append(args, "--format", "yaml")
	log.Println("gcloud", strings.Join(args, " "))
	cmd := exec.Command("gcloud", args...)
	// pass user's PATH env variable, expected gcloud executable can be found in PATH
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("failed to start cmd: %v, \n output: %v", err, string(output))
		return string(output), err
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("cmd returned error: %v, \n output: %v", err, string(output))
		return string(output), err
	}

	return string(output), err
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

// Create deployment based on passed Deployment object passed into it.
// Initialize Deployment with Outputs map in case of successful creation and error otherwise.
func CreateOrUpdate(action string, deployment *Deployment) (string, error) {
	if action != ActionCreate && action != ActionUpdate {
		log.Fatalf("action %s not in [%s,%s] for deployment: %v", action, ActionCreate, ActionUpdate, deployment)
	}

	args := []string{
		"deployment-manager",
		"deployments",
		action,
		deployment.config.Name,
		"--config",
		deployment.configFile,
		"--project",
		deployment.config.Project,
	}
	output, err := runGCloud(args...)
	if err != nil {
		log.Printf("failed to %s deployment: %v, error: %v", action, deployment, err)
		return output, err
	}
	outputs, err := GetOutputs(deployment.config.Name, deployment.config.Project)
	if err != nil {
		log.Printf("on %s action, failed to get outputs for deployment: %v, error: %v", action, deployment, err)
		return output, err
	}
	deployment.Outputs = outputs
	return output, nil
}

// Delete function removed Deployment passed into it as parameter.
func Delete(deployment *Deployment) (string, error) {
	args := []string{
		"deployment-manager",
		"deployments",
		"delete",
		deployment.config.Name,
		"--project",
		deployment.config.Project,
		"-q",
	}
	output, err := runGCloud(args...)
	if err != nil {
		log.Printf("failed to get deployment manifest: %v", err)
		return output, err
	}
	return output, nil
}

func GetStatus(deployment *Deployment) (Status, error) {
	args := []string{
		"deployment-manager",
		"deployments",
		"describe",
		deployment.config.Name,
		"--project",
		deployment.config.Project,
	}
	response, err := runGCloud(args...)
	if err != nil {
		if strings.Contains(response, "code=404") {
			return NotFound, nil
		} else {
			log.Printf("failed to get status for deployment: %s, \n error: %v", deployment.config.FullName(), err)
			return Error, err
		}
	}

	description := &struct {
		Deployment struct {
			Name      string
			Operation struct {
				Status        string
				OperationType string
			}
		}
	}{}

	err = yaml.Unmarshal([]byte(response), description)
	if err != nil {
		log.Printf("error unmarshall response: %s,\n deployment: %v \n error: %v", response, deployment, err)
		return Error, err
	}

	status := description.Deployment.Operation.Status

	switch status {
	case "DONE":
		return Done, nil
	case "RUNNING":
		return Running, nil
	case "PENDING":
		return Pending, nil
	default:
		return Error, errors.New(fmt.Sprintf("Unknown status %s, for deployment %s",
			deployment.config.FullName(), status))
	}
}

func parseOutputs(data string) (map[string]string, error) {
	describe, err := unmarshal(data)
	if err != nil {
		log.Println("error parsing deployment outputs")
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
		log.Println("error parsing deployment outputs layout section")
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
