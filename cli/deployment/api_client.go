package deployment

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// Status represent Deployment status in enum/numerical format
type Status int

const (
	// Done Status represent successfully finished Deployment state.
	Done Status = 0
	// Pending Status means Deployment in "PENDING" state, means resource creation process not started yet.
	Pending Status = 1
	// Pending Status means Deployment in "RUNNING" state, means resource creation process started.
	Running Status = 2
	// NotFound Status means Deployment is not exists in any state.
	NotFound Status = 3
	// Error Status means Deployment is failed or deployment describe operation itself completed with error.
	Error Status = 4
)

const (
	// ApplyAction passed as "action" parameter value to cmd.execute(action string, ...) to apply Deployment created or updated in preview mode.
	ActionApply string = "apply"
	// ApplyAction passed as "action" parameter value to cmd.execute(action string, ...) to delete Deployment.
	ActionDelete string = "delete"
	// ApplyAction passed as "action" parameter value to cmd.execute(action string, ...) to update Deployment.
	ActionCreate string = "create"
	// ApplyAction passed as "action" parameter value to cmd.execute(action string, ...) to update Deployment.
	ActionUpdate string = "update"
)

type DeploymentDescriptionResource struct {
	Name            string
	Type            string
	Properties      string
	FinalProperties string `yaml:",omitempty"`
	Update          struct {
		Properties      string
		FinalProperties string `yaml:",omitempty"`
		State           string
	} `yaml:",omitempty"`
}

type DeploymentDescription struct {
	Deployment struct {
		Name      string
		Operation struct {
			OperationType string
			Status        string
		}
	}
	Resources []DeploymentDescriptionResource
}

// Resource struct used to parse deployment manifest yaml received with 'gcloud deployment-manager manifests describe' command.
type Resources struct {
	Name    string
	Outputs []struct {
		Value interface{} `yaml:"finalValue"`
		Name  string
	}
	Resources []Resources
}

var actions = []string{ActionApply, ActionDelete, ActionCreate, ActionUpdate}

// String returns human readable string representation of Deployment Status.
func (s Status) String() string {
	return [...]string{"DONE", "PENDING", "RUNNING", "NOT_FOUND", "ERROR"}[s]
}

// The RunGCloud function runs the gcloud tool with the specified arguments. It is implemented
// as a variable so that it can be mocked in tests of its exported consumers.
var RunGCloud = func(args ...string) (result string, err error) {
	log.Println("gcloud", strings.Join(args, " "))
	cmd := exec.Command("gcloud", args...)
	// pass user's PATH env variable, expected gcloud executable can be found in PATH
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))

	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), err
	}

	return string(output), err
}

// GetOutputs retrieves existing Deployment outputs using gcloud and store result in map[string]interface{}
// where "resourceName.propertyName" is key, and value is string (in case of flat value) or JSON object.
func GetOutputs(project string, name string) (map[string]interface{}, error) {
	output, err := RunGCloud("deployment-manager", "manifests", "describe", "--deployment", name, "--project", project, "--format", "yaml")
	if err != nil {
		log.Printf("failed to describe deployment manifest for deployment: %s.%s, error: %v, output: %s", project, name, err, output)
		return nil, err
	}
	return parseOutputs(output)
}

// GCloudDefaultProjectID returns the default project id taken from local gcloud configuration.
func GCloudDefaultProjectID() (string, error) {
	data, err := RunGCloud("config", "list", "--format", "yaml")
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

// Create creates deployment based on passed Deployment object passed into it.
// Create initialize passed Deployment object with Outputs map in case of successful creation and return error otherwise.
// preview parameter define if deployment should be created in 'Preview' mode.
// Function returns gcloud cli raw output for debug purposes both in case of success and error.
func Create(deployment *Deployment, preview bool) (string, error) {
	return createOrUpdate(ActionCreate, deployment, preview)
}

// Update updates deployment based on passed Deployment object passed into it.
// Update initialize passed Deployment object with Outputs map in case of successful creation and return error otherwise.
// preview parameter define if deployment should be updated in 'Preview' mode.
// Function returns gcloud cli raw output for debug purposes both in case of success and error.
func Update(deployment *Deployment, preview bool) (string, error) {
	return createOrUpdate(ActionUpdate, deployment, preview)
}

func createOrUpdate(action string, deployment *Deployment, preview bool) (string, error) {
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

	output, err := RunGCloud(args...)
	if err != nil {
		log.Printf("failed to %s deployment: %v, error: %v, output: %s", action, deployment, err, output)
		return output, err
	}

	if !preview {
		outputs, err := GetOutputs(deployment.config.GetProject(), deployment.config.Name)
		if err != nil {
			log.Printf("on %s action, failed to get outputs for deployment: %v, error: %v, output: %s", action, deployment, err, output)
			return output, err
		}
		deployment.Outputs = outputs
	}

	return output, nil
}

// CancelPreview cancels update/create/delete action, created with review flag.
// Function uses gcloud deployments cancel-preview command.
// In case of cancellation of preview of create action, required deployment and run Delete() after CancelPreview() for cleanup.
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
	output, err := RunGCloud(args...)
	if err != nil {
		log.Printf("failed to cancel preview deployment: %v, error: %v, output: %s", deployment, err, output)
		return output, err
	}
	return output, nil
}

// ApplyPreview function apply changes made before with --preview flag
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
	output, err := RunGCloud(args...)
	if err != nil {
		log.Printf("failed to apply preview for deployment: %v, error: %v, output: %s", deployment, err, output)
		return output, err
	}
	return output, nil
}

// Delete function removed Deployment passed into it as parameter.
// Boolean preview param define if changes have to be previewed before apply.
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
	output, err := RunGCloud(args...)
	if err != nil {
		log.Printf("failed to delete deployment: %v, error: %v, output: %s", deployment, err, output)
		return output, err
	}
	return output, nil
}

func GetDeploymentDescription(name string, project string) (*DeploymentDescription, error) {
	args := []string{
		"deployment-manager",
		"deployments",
		"describe",
		name,
		"--project",
		project,
		"--format", "yaml",
	}
	response, err := RunGCloud(args...)
	if err != nil {
		return nil, err
	}

	description := &DeploymentDescription{}

	err = yaml.Unmarshal([]byte(response), description)
	if err != nil {
		log.Printf("error unmarshall response: %s,\n deployment: %v \n error: %v", response, name, err)
		return nil, err
	}

	return description, nil
}

// GetStatus retrieves Deployment status using gcloud cli, see deployment.Status type for details.
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
	response, err := RunGCloud(args...)
	if err != nil {
		if strings.Contains(response, "code=404") {
			return NotFound, nil
		} else {
			return Error, err
		}
	}

	description := &DeploymentDescription{}

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
		return Error, fmt.Errorf("Unknown status %s, for deployment %s",
			deployment.config.FullName(), status)
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
		Resources []Resources
	}{}
	err = yaml.Unmarshal([]byte(layoutData), res)
	if err != nil {
		log.Println("error parsing deployment outputs layout section")
		return nil, err
	}

	result := make(map[string]interface{})

	resources := flattenResources(res.Resources)
	for _, resource := range resources {
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

// flattenResources iterate over passed slice of Resources object, iterates over all sub-resources recursively and add all
// resources to result array. In simple worlds flattenResources extracts all resouces and sub-resources with non empty Outputs field.
func flattenResources(source []Resources) []Resources {
	var result []Resources
	for _, resource := range source {
		if len(resource.Outputs) > 0 {
			result = append(result, resource)
		}
		if len(resource.Resources) > 0 {
			result = append(result, flattenResources(resource.Resources)...)
		}
	}
	return result
}
