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
	Pending  Status = 111
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

// function exposed to variable in order to mock it inside api_client_test.go
var runGCloud = func(args ...string) (result string, err error) {
	log.Println("gcloud", strings.Join(args, " "))
	cmd := exec.Command("gcloud", args...)
	// pass user's PATH env variable, expected gcloud executable can be found in PATH
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("failed to start cmd: %v, output: %v", err, string(output))
		return string(output), err
	}

	return string(output), err
}

// GetOutputs execute deployment-manager manifest describe call with gcloud tool and parse returned
// resources.output yaml section. Returns map where "resourceName.propertyName" is key
func GetOutputs(project string, name string) (map[string]interface{}, error) {
	data, err := runGCloud("deployment-manager", "manifests", "describe", "--deployment", name, "--project", project, "--format", "yaml")
	if err != nil {
		log.Printf("Failed to get deployment manifest: %v", err)
		return nil, err
	}
	return parseOutputs(data)
}

/*
returns project id taken from local gcloud configuration
*/
func GCloudDefaultProjectID() (string, error) {
	data, err := runGCloud("config", "list", "--format", "yaml")
	if err != nil {
		return "", err
	}
	out := struct {
		Core struct {
			Project string
		}
	}{}
	err = yaml.Unmarshal([]byte(data), &out)
	if err != nil {
		return "", err
	}
	return out.Core.Project, nil
}

// Create deployment based on passed Deployment object
// Returns Deployment with Outputs map in case of successful creation and error otherwise
func CreateOrUpdate(action string, deployment *Deployment, preview bool) (string, error) {
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
		deployment.config.GetProject(),
	}

	if preview {
		args = append(args, "--preview")
	}

	output, err := runGCloud(args...)
	if err != nil {
		log.Printf("failed to %s deployment: %v, error: %v", action, deployment, err)
		return output, err
	}

	if !preview {
		outputs, err := GetOutputs(deployment.config.GetProject(), deployment.config.Name)
		if err != nil {
			log.Printf("on %s action, failed to get outputs for deployment: %v, error: %v", action, deployment, err)
			return output, err
		}
		deployment.Outputs = outputs
	}
	return output, nil
}

/*
cancel update/create/delete preview with gcloud deployments cancel-preview command,
in case of cancel preview of create action, need to clean deployment and run Delete() after CancelPreview()
*/
func CancelPreview(deployment *Deployment) (string, error) {
	args := []string{
		"deployment-manager",
		"deployments",
		"cancel-preview",
		deployment.config.Name,
		"--project",
		deployment.config.GetProject(),
		"-q",
	}
	output, err := runGCloud(args...)
	if err != nil {
		log.Printf("failed to cancel preview, error: %v", err)
		return output, err
	}
	return output, nil
}

/*
Applies changes made before with --preview flag
*/
func ApplyPreview(deployment *Deployment) (string, error) {
	args := []string{
		"deployment-manager",
		"deployments",
		"update",
		deployment.config.Name,
		"--project",
		deployment.config.GetProject(),
		"-q",
	}
	output, err := runGCloud(args...)
	if err != nil {
		log.Printf("failed to apply preview, error: %v", err)
		return output, err
	}
	return output, nil
}

func Delete(deployment *Deployment, preview bool) (string, error) {
	args := []string{
		"deployment-manager",
		"deployments",
		"delete",
		deployment.config.Name,
		"--project",
		deployment.config.GetProject(),
		"-q",
	}
	if preview {
		args = append(args, "--preview")
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
		deployment.config.GetProject(),
		"--format", "yaml",
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
		log.Printf("Error unmarshall response: %s,\n deployment: %v \n error: %v", response, deployment, err)
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

func parseOutputs(data string) (map[string]interface{}, error) {
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

	result := make(map[string]interface{})
	for _, resource := range res.Resources {
		for _, output := range resource.Outputs {
			key := resource.Name + "." + output.Name
			value := output.Value
			result[key] = value
		}
	}

	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}
