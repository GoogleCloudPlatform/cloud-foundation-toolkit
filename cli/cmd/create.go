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
	cmd.Printf("Create deployment, configs %v, project %s\n", args, ProjectFlag)
	configs := loadConfigs(args)
	//graph := deployment.NewDependencyGraph(configs)
	// TODO fix graph sort
	// ordered, err := graph.Order()
	// if err != nil {
	//	log.Fatal("Error during creating deployment dependencies graph", err)
	// }

	ordered := configs

	outputs := make(map[string]map[string]string)
	for _, config := range ordered {
		dep := deployment.NewDeployment(config, outputs)
		log.Println("Start creating deployment", dep.String())
		result, err := deployment.Create(dep)
		if err != nil {
			log.Fatalf("Error during creating deployment %v, \n %v", dep, err)
		}
		outputs[result.ID()] = result.Outputs
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
