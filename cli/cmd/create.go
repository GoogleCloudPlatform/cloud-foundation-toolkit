package cmd

import (
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var ProjectFlag string

func init() {
	createCmd.PersistentFlags().StringVarP(&ProjectFlag, "project", "p", "", "project name")
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create deployment(s)",
	Long:  `Create deployment(s)`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Printf("Create deployment, configs %v, project %s", args, ProjectFlag)
	configs := loadConfigs(args)
	graph := deployment.NewDependencyGraph(configs)
	ordered, err := graph.Order()
	if err != nil {
		log.Fatal("Error during creating deployment dependencies graph", err)
	}
	for _, config := range ordered {
		dep := deployment.NewDeployment(config)
		_, err = deployment.Create(dep)
		if err != nil {
			log.Fatal("Error during creating deployment %v, \n %v", dep, err)
		}
	}
}

func loadConfigs(args []string) []deployment.Config {
	result := make([]deployment.Config, len(args))
	for i, path := range args {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("Error loading file", path)
		}
		result[i] = *deployment.NewConfig(string(data), path)
	}
	return result
}
