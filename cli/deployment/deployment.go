package deployment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

// The Deployment type represents a real GCP Deployment entity that is either already created, or yet-to-be created.
type Deployment struct {
	// Outputs map contains deployment outputs values in form resourceName.proeprtyName: value, map filled with data
	// after deployment update/create operation. Value could be either plain string or complex object.
	Outputs map[string]interface{}
	// config object store config state parsed from config YAML as it is, no modification, cross-deployment reference values substitution etc
	config Config
	// configFile field point to YAML file, generated from Config object with all cross-deployment references
	// overwritten with actual values.
	configFile string
}

// NewDeployment creates a new Deployment object, overriding all outward references.
// In effect, this means all deployment dependencies must exist in the GCP project.
// output parameter is map of maps where key is Deployment full name (project_name.deployment_name), and value is map
// of corresponding Deployment outputs properties.
func NewDeployment(config Config, outputs map[string]map[string]interface{}, processConfig bool) *Deployment {
	file, err := ioutil.TempFile("", config.Name)
	defer func() {
		er := file.Close()
		if er != nil {
			log.Printf("close temp config file error : %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("Error creating temp file for deployment: %s, error: %v", config.FullName(), err)
	}
	var data []byte
	if processConfig {
		// don't need to process config, fix import paths, replace out refs in case of Delete command,
		// only deployment name and project required for operation
		data = replaceOutRefs(config, outputs)
		_, err = file.Write(data)
		if err != nil {
			log.Fatalf("Error wirte to file: %s, error: %v", file.Name(), err)
		}
	}
	return &Deployment{
		config:     config,
		configFile: file.Name(),
	}
}

// String implements fmt.Stringer for the Deployment type.
func (d Deployment) String() string {
	return fmt.Sprintf("Deployment[name=%s, config=%s]", d.config.FullName(), d.configFile)
}

// FullName function is the same as deployment.Config.FullName(), can be used in map[string]Deployment as a key.
func (d Deployment) FullName() string {
	return d.config.FullName()
}

// Executes executes create/update/delete actions and returns Deployment status.
func (d *Deployment) Execute(action string, preview bool) (output string, error error) {
	if sort.SearchStrings(actions, action) == len(actions) {
		log.Fatalf("action: %s not in %v for deployment: %v", actions, actions, d)
	}

	switch action {
	case ActionCreate:
		return Create(d, preview)
	case ActionUpdate:
		return Update(d, preview)
	case ActionDelete:
		return Delete(d, preview)
	case ActionApply:
		status, err := GetStatus(d)
		if err != nil {
			log.Printf("Apply action for deployment: %s, break error: %v", d, err)
		}
		switch status {
		case Done:
			log.Printf("Deployment %v exists, run Update()", d)
			return Update(d, preview)
		case NotFound:
			log.Printf("Deployment %v does not exists, run Create()", d)
			return Create(d, preview)
		case Pending:
			log.Printf("Deployment %v is in pending state, break", d)
			return "", fmt.Errorf("deployment %v is in PENDING state", d)
		case Error:
			message := fmt.Sprintf("Could not get state of deployment: %v", d)
			log.Print(message)
			return "", errors.New(message)
		}
		return "", fmt.Errorf("error during Apply command for deployment: %v", d)
	default:
		log.Fatalf("unrecognized action %s", action)
	}
	return "", nil
}

func replaceOutRefs(config Config, outputs map[string]map[string]interface{}) []byte {
	data, err := config.YAML(outputs)
	if err != nil {
		log.Fatalf("error while parsing yaml for config: %s, error: %v", config.FullName(), err)
	}
	return data
}
