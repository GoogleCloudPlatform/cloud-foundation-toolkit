package deployment

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"log"
	"os/exec"
)

// during testing it will be replaced with mock
var execFunction = func(args ...string) (result string, err error) {
	cmd := exec.Command("/home/kopachevsky/google-cloud-sdk/bin/gcloud", append(args, "--format", "yaml")...)
	//cmd.Dir = entryPath
	outBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}
	cmd.Stdout = outBuff
	cmd.Stderr = errBuff

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

func GetOutputs(name string, project string) (map[string]string, error) {
	data, err := execFunction("deployment-manager", "manifests", "describe", "--deployment", name, "--project", project)
	if err != nil {
		log.Print("Failed to get deployment manifest", err)
		return nil, err
	}
	return ParseOutputs(data)
}

func Create(deployment *Deployment) (*Deployment, error) {
	_, err := execFunction("deployment-manager", "deployments", "create", deployment.config.Name, "--config",
		deployment.ConfigFile(), "--project", deployment.config.Project)
	if err != nil {
		log.Print("Failed to create deployment", err)
		return nil, err
	}
	outputs, err := GetOutputs(deployment.config.Name, deployment.config.Project)
	if err != nil {
		log.Print("Failed to get deployment outputs", err)
		return nil, err
	}
	deployment.outputs = outputs
	return deployment, nil
}

func ParseOutputs(data string) (map[string]string, error) {

	type resources struct {
		Resources []struct {
			Name    string
			Outputs []struct {
				Value interface{}
				Name  string
			}
		}
	}

	describe, err := unmarshal(data)
	if err != nil {
		log.Println("Error unmarshal deployment outputs")
		return nil, err
	}

	layoutData := describe["layout"].(string)

	res := &resources{}
	err = yaml.Unmarshal([]byte(layoutData), res)
	if err != nil {
		log.Println("Error unmarshal deployment outputs layout data")
		return nil, err
	}

	result := make(map[string]string)
	for _, resource := range res.Resources {
		for _, output := range resource.Outputs {
			switch value := output.Value.(type) {
			case string:
				result[resource.Name+"."+output.Name] = value
			case map[interface{}]interface{}:
				log.Println(resource.Name + "." + output.Name + " is map")
			}
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}
