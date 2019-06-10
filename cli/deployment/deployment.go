package deployment

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// The Deployment type represents a real GCP Deployment entity that is either already created, or yet-to-be created.
type Deployment struct {
	// Outputs map contains deployment outputs values in form resourceName.proeprtyName: value, map filled with data
	// after deployment update/create operation
	Outputs map[string]string
	// config object store config state parsed from config YAML as it is, no modification, cross-deployment reference values substitution etc
	config Config
	// configFile field point to YAML file, generated from Config object with all cross-deployment references
	// overwritten with actual values.
	configFile string
}

// NewDeployment creates a new Deployment object, overriding all outward references.
// In effect, this means all deployment dependencies must exist in the GCP project.
func NewDeployment(config Config, outputs map[string]map[string]string) *Deployment {
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

	data := replaceOutRefs(config, outputs)

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Error wirte to file: %s, error: %v", file.Name(), err)
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

func (d *Deployment) Execute(action string) (output string, error error) {
	if sort.SearchStrings(actions, action) == len(actions) {
		log.Fatalf("action: %s not in %v for deployment: %v", actions, actions, d)
	}

	if action == ActionCreate || action == ActionUpdate {
		return CreateOrUpdate(action, d)
	} else if action == ActionDelete {
		return Delete(d)
	} else {
		status, err := GetStatus(d)
		if err != nil {
			log.Printf("Apply action for deployment: %s, break error: %v", d, err)
		}
		switch status {
		case Done:
			log.Printf("Deployment %v exists, run Update()", d)
			return CreateOrUpdate(ActionUpdate, d)
		case NotFound:
			log.Printf("Deployment %v does not exists, run Create()", d)
			return CreateOrUpdate(ActionCreate, d)
		case Pending:
			log.Printf("Deployment %v is in pending state, break", d)
			return "", errors.New(fmt.Sprintf("Deployment %v is in PENDING state", d))
		case Error:
			message := fmt.Sprintf("Could not get state of deployment: %v", d)
			log.Print(message)
			return "", errors.New(message)
		}
		return "", errors.New(fmt.Sprintf("Error during Apply command for deployment: %v", d))
	}
}

func replaceOutRefs(config Config, outputs map[string]map[string]string) []byte {
	data, err := config.YAML()
	if err != nil {
		log.Fatalf("error while parsing yaml for config: %s, error: %v", config.FullName(), err)
	}
	refs := config.findAllOutRefs()
	for _, ref := range refs {
		fullName, resource, property := parseOutRef(ref)
		outputsMap := outputs[fullName]
		if outputsMap == nil {
			arr := strings.Split(fullName, ".")
			outputsMap, err := GetOutputs(arr[0], arr[1])
			if err != nil {
				log.Fatalf("Erorr getting outputs for deployment: %s, error: %v", fullName, err)
			}
			outputs[fullName] = outputsMap
		}
		key := resource + "." + property
		value := outputsMap[key]
		fullRef := fmt.Sprintf("$(out.%s)", ref)
		if len(value) == 0 {
			log.Fatalf("Could not resolve reference: %s", fullRef)
		}
		data = bytes.ReplaceAll(data, []byte(fullRef), []byte(value))
	}
	return data
}
