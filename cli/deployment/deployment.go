package deployment

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Deployment struct {
	outputs map[string]string
	config  Config
	// tmp config file location
	configFile string
}

func (d Deployment) ConfigFile() string {
	return d.configFile
}

func NewDeployment(config Config, outputs map[string](map[string]string)) *Deployment {
	file, err := ioutil.TempFile("dm", config.Name)
	if err != nil {
		log.Fatal(err)
	}
	path, _ := filepath.Abs(filepath.Dir(file.Name()))
	ioutil.WriteFile(path, []byte(config.data), os.ModeTemporary)

	return &Deployment{
		config:     config,
		configFile: path,
	}
}

func (d Deployment) ReplaceOutRefs(data []byte, outputs map[string](map[string]string)) {
	refs := d.config.findAllOutRefs()
	for _, ref := range refs {
		project, deployment, _, _ := parseOutRef(ref)
		outputsMap := outputs[project+"."+deployment]
		if outputsMap == nil {
			outputsMap, err := GetOutputs(deployment, project)
			if err != nil {
				log.Fatal(err)
			}
			outputs[project+"."+deployment] = outputsMap
		}
		value := outputsMap[ref]
		fullRef := fmt.Sprint("${out.%s}", ref)
		if len(value) == 0 {
			log.Fatal("Could not resolve reference ", fullRef)
		}
		bytes.ReplaceAll(data, []byte(fullRef), []byte(value))
	}
}

func (d Deployment) String() string {
	return fmt.Sprintf("Deployment[name=%s, config=%s]", d.config.String(), d.configFile)
}
