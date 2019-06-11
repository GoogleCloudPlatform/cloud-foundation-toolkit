package cmd

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

// --project flag value will be mapped during app initialization
var projectFlag string

// common code for create/update/apply and delete actions
func execute(action string, cmd *cobra.Command, args []string) {
	cmd.Printf("%s deployment(s), configs %v, project %s\n", action, args, projectFlag)
	configs := loadConfigs(args)
	ordered, err := deployment.Order(configs)
	if err != nil {
		log.Fatalf("Error ordering deployments in dependency order: %v", err)
	}
	isDelete := action == deployment.ActionDelete
	if isDelete {
		// reverse order, dependent goes first for deletion
		for i := len(ordered)/2 - 1; i >= 0; i-- {
			opp := len(ordered) - 1 - i
			ordered[i], ordered[opp] = ordered[opp], ordered[i]
		}
	}
	log.Printf("Ordered dependencies: %v", ordered)

	outputs := make(map[string]map[string]interface{})
	for i, config := range ordered {
		dep := deployment.NewDeployment(config, outputs, !isDelete)
		cmd.Printf("---------- Stage %d ----------\n", i)
		output, err := dep.Execute(action)
		cmd.Print(output)
		if err != nil {
			log.Fatalf("Error %s deployment %v, \n %v", action, dep, err)
		}
		// after create/update/apply actions - put deployment outputs to global map fro cross-deployment
		// reference variables substitutions, after delete - remove deployment outputs from map, to avoid its usage
		if action != deployment.ActionDelete {
			outputs[dep.FullName()] = dep.Outputs
		} else {
			delete(outputs, dep.FullName())
		}
	}
}

func loadConfigs(args []string) map[string]deployment.Config {
	result := make(map[string]deployment.Config, len(args))
	for _, path := range args {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Error loading file: %s", path)
		}
		config := deployment.NewConfig(string(data), path)
		result[config.FullName()] = config
	}
	return result
}
