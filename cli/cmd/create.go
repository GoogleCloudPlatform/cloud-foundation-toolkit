package cmd

import (
	"log"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

var projectFlag string

func init() {
	createCmd.PersistentFlags().StringVarP(&projectFlag, "project", "p", "", "project name")
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create deployment(s)",
	Long:  `Create deployment(s)`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Printf("Create deployment, configs %v, project %s\n", args, projectFlag)
	configs := loadConfigs(args)
	ordered, err := deployment.Order(configs)
	if err != nil {
		log.Fatalf("Error ordering deployments in dependency order: %v", err)
	}

	log.Printf("Ordered dependencies: %v", ordered)

	outputs := make(map[string]map[string]string)
	for _, config := range ordered {
		dep := deployment.NewDeployment(config, outputs)
		log.Println("Start creating deployment", dep.String())
		result, err := deployment.Create(dep)
		if err != nil {
			log.Fatalf("Error during creating deployment %v, \n %v", dep, err)
		}
		outputs[result.FullName()] = result.Outputs
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
