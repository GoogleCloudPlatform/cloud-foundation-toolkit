package deployment

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os/exec"
)

type Deployment struct {
	name    string
	outputs map[string]string
}

func unmarshal(data string) (map[string]interface{}, error) {
	my := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(data), my)
	if err != nil {
		fmt.Printf("error %+v", err)
		return nil, err
	}
	return my, nil
}

func GCloud(args ...string) (result string, err error) {
	cmd := exec.Command("/home/kopachevsky/google-cloud-sdk/bin/gcloud", append(args, "--format", "yaml")...)
	//cmd.Dir = entryPath
	outBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}
	cmd.Stdout = outBuff
	cmd.Stderr = errBuff

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
		return "", err
	}

	log.Println("Doing other stuff...")

	if err := cmd.Wait(); err != nil {
		log.Printf("Cmd returned error: %v", err)
		return errBuff.String(), err
	}

	return outBuff.String(), err
}

func GetOutputs(data string) (map[string]string, error) {

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
