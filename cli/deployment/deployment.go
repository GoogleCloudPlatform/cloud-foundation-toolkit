package deployment

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

type Deployment struct {
	Outputs    map[string]string
	config     Config
	configFile string
}

func (d Deployment) ConfigFile() string {
	return d.configFile
}

func NewDeployment(config Config, outputs map[string]map[string]string) *Deployment {
	file, err := ioutil.TempFile("", config.Name)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	data := replaceOutRefs(config, outputs)

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	return &Deployment{
		config:     config,
		configFile: file.Name(),
	}
}

func replaceOutRefs(config Config, outputs map[string]map[string]string) []byte {
	data, err := config.Yaml()
	if err != nil {
		panic(err)
	}
	refs := config.findAllOutRefs()
	for _, ref := range refs {
		project, deployment, resource, property := parseOutRef(ref)
		outputsMap := outputs[project+"."+deployment]
		if outputsMap == nil {
			outputsMap, err := GetOutputs(deployment, project)
			if err != nil {
				log.Fatal(err)
			}
			outputs[project+"."+deployment] = outputsMap
		}
		key := resource + "." + property
		value := outputsMap[key]
		fullRef := fmt.Sprintf("$(out.%s)", ref)
		if len(value) == 0 {
			log.Fatal("Could not resolve reference ", fullRef)
		}
		data = bytes.ReplaceAll(data, []byte(fullRef), []byte(value))
	}
	return data
}

func (d Deployment) String() string {
	return fmt.Sprintf("Deployment[name=%s, config=%s]", d.config.String(), d.configFile)
}

func (d Deployment) ID() string {
	return d.config.String()
}
