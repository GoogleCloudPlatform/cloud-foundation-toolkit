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
var previewFlag bool = false

var supportedExt = []string{"*.yaml", "*.yml", "*.jinja"}

func initCommon(command *cobra.Command) {
	command.PersistentFlags().StringVarP(&projectFlag, "project", "p", "", "project id")
	command.PersistentFlags().BoolVar(&previewFlag, "preview", false, "preview before apply changes")
	rootCmd.AddCommand(command)
}

// common code for create/update/apply and delete actions
func execute(action string, cmd *cobra.Command, args []string) {
	setDefaultProjectID()
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
		output, err := dep.Execute(action, previewFlag)
		cmd.Print(output)
		if err != nil {
			if action == deployment.ActionDelete {
				status, _ := deployment.GetStatus(dep)
				if status == deployment.NotFound {
					// for Delete action, Deployment might not exists, in this case just skip
					log.Printf("Deployment %v does not exists, skip deletion\n", dep)
					continue
				}
			}
			log.Fatalf("Error %s deployment: %v, erro: %v", action, dep, err)
		}
		if previewFlag {
			choise := deployment.GetUserInput("Update(u), Skip (s), or Abort(a) Deployment?", []string{"u", "s", "a"}, os.Stdin)
			switch choise {
			case "u":
				output, err := deployment.ApplyPreview(dep)
				cmd.Print(output)
				if err != nil {
					log.Fatalf("error %s deployment: %v, error: %v", action, dep, err)
				}
			case "s":
				output, err = deployment.CancelPreview(dep)
				if err != nil {
					log.Fatalf("error cancel-preuvew action for deployment: %v, error: %v", dep, err)
				}
				log.Printf("canceled %s action for deployment: %v", action, dep)
				if action == deployment.ActionCreate {
					log.Printf("delete canceled creation for deployment: %v", dep)
					deployment.Delete(dep, false)
				}
			case "a":
				log.Print("Aborting deployment run!")
				os.Exit(0)
			}
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

/**
 listConfigs search for config files according rules described in CLL usage section of following doc:
https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/dm/docs/userguide.md#syntax
listConfigs returns map[fileName: fileContent] for files, and list of strings for yamls passed as string parameters to cli
*/
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
					errs[arg] = errors.New(fmt.Sprintf("no %s files found in directory: %s", strings.Join(supportedExt, ", "), arg))
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

/*
listConfigs accept list of config parameters (file/directory paths, glob pattern, yaml strings)
search all possible files in case of directory/glob patterns with listConfigs function and create
Config objects from loaded data
*/
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
		for _, data := range yamls {
			config := deployment.NewConfig(data, "")
			result[config.FullName()] = config
		}
	}

	if len(result) == 0 {
		log.Fatal("no configs provided")
	}
	return result
}

/*
set deployment.DefaultProjectID variable by search following options:
The --project command-line option.
The CLOUD_FOUNDATION_PROJECT_ID environment variable.
The "default project" configured with the GCP SDK.
*/
func setDefaultProjectID() {
	if len(projectFlag) > 0 {
		deployment.DefaultProjectID = projectFlag
	} else if env := os.Getenv("CLOUD_FOUNDATION_PROJECT_ID"); len(env) > 0 {
		deployment.DefaultProjectID = env
	} else {
		gcloudDefault, err := deployment.GCloudDefaultProjectID()
		if err != nil {
			log.Fatalf("error getting gcloud default project: %v", err)
		}
		if len(gcloudDefault) == 0 {
			log.Fatalf("can't get project id from --project arg, CLOUD_FOUNDATION_PROJECT_ID env variable and gcloud default")
		}
		deployment.DefaultProjectID = gcloudDefault
	}
}
