package deployment

import (
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

func NewDeployment(config Config) *Deployment {

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

func (d Deployment) String() string {
	return fmt.Sprintf("Deployment[name=%s, config=%s]", d.config.String(), d.configFile)
}
