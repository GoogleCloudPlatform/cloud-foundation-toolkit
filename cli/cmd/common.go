package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
)

// --project flag value will be mapped during app initialization
var projectFlag string

var supportedExt = []string{"*.yaml", "*.yml", "*.jinja"}

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

func listConfigs(args []string, errs map[string]error) (map[string]string, []string) {
	resFiles := map[string]string{}
	var resYamls []string
	for _, arg := range args {
		stat, err := os.Stat(arg)
		if err == nil {
			if stat.IsDir() {
				configs := map[string]string{}
				for _, ext := range supportedExt {
					glob := path.Clean(arg) + "/" + ext
					files, _ := listConfigs([]string{glob}, errs)
					configs = deployment.AppendMap(configs, files)
				}
				resFiles = deployment.AppendMap(resFiles, configs)
				if len(configs) == 0 {
					errs[arg] = errors.New(fmt.Sprintf("No %s files found in directory: %s", strings.Join(supportedExt, ", "), arg))
				}
			} else {
				data, err := ioutil.ReadFile(arg)
				if err != nil {
					log.Fatalf("file %s read error: %v", arg, err)
				}
				resFiles[arg] = string(data)
			}
		} else if os.IsNotExist(err) {
			if deployment.IsYAML(arg) {
				resYamls = append(resYamls, arg)
			} else {
				// check Glob
				maches, err := filepath.Glob(arg)
				if err != nil {
					errs[arg] = errors.New(fmt.Sprintf("Error during search files for config: %s, %v", arg, err))
				} else {
					if len(maches) > 0 {
						files, _ := listConfigs(maches, errs)
						resFiles = deployment.AppendMap(resFiles, files)
					} else {
						errs[arg] = errors.New(fmt.Sprintf("No file(s) exists or valid yaml for config param: %s", arg))
					}
				}
			}
		} else {
			log.Fatalf("file %s stat error: %v", arg, err)
		}
	}
	return resFiles, resYamls
}

func loadConfigs(args []string) map[string]deployment.Config {
	result := map[string]deployment.Config{}
	errs := map[string]error{}
	files, yamls := listConfigs(args, errs)

	// check errors
	for _, entry := range args {
		if err, ok := errs[entry]; ok {
			log.Fatal(err)
		}
	}

	for file, data := range files {
		config := deployment.NewConfig(data, file)
		result[config.FullName()] = config
	}

	if len(yamls) > 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not get current folder path: %v", err)
		}
		for _, data := range yamls {
			config := deployment.NewConfig(data, dir)
			result[config.FullName()] = config
		}
	}

	if len(result) == 0 {
		log.Fatal("no configs provided")
	}

	return result
}
