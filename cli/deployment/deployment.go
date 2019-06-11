package deployment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

// Deployment object represent real GCP Deployment entity that already created or have to be created
// configFile field point to yaml file, generated from Config object with all cross-deployment references
// overwritten with actual values
type Deployment struct {
	Outputs    map[string]interface{}
	config     Config
	configFile string
}

// NewDeployment creates Deployment object and override all out refs, this means all
// deployments it depends to should exists in GCP project
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

// String representation of Deployment instance
func (d Deployment) String() string {
	return fmt.Sprintf("Deployment[name=%s, config=%s]", d.config.FullName(), d.configFile)
}

// same as deployment.Config.FullName(), can be used in map[string]Deployment as a key
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

func replaceOutRefs(config Config, outputs map[string]map[string]interface{}) []byte {
	data, err := config.YAML(outputs)
	if err != nil {
		log.Fatalf("error while parsing yaml for config: %s, error: %v", config.FullName(), err)
	}
	return data
}
