package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation"
)

// --project flag value will be mapped during app initialization
var projectFlag string
var policyPathFlag string
var validateFlag bool = false
var previewFlag bool = false
var showStagesFlag bool = false
var formatFlag string

var supportedExt = []string{"*.yaml", "*.yml", "*.jinja"}

func initPolicyPathFlag(command *cobra.Command) {
	command.PersistentFlags().StringVar(&policyPathFlag, "policy-path", "", "Policy path")
}

func initProjectFlag(command *cobra.Command) {
	command.PersistentFlags().StringVarP(&projectFlag, "project", "p", "", "project id")
}

func initValidateFlags(command *cobra.Command) {
	initPolicyPathFlag(command)
	command.PersistentFlags().BoolVar(&validateFlag, "validate", false, "validate deployment")
}

func initCommon(command *cobra.Command) {
	initProjectFlag(command)
	command.PersistentFlags().BoolVar(&previewFlag, "preview", false, "preview before apply changes")
	command.PersistentFlags().BoolVar(&showStagesFlag, "show-stages", false, "print deployment stages")
	command.PersistentFlags().StringVar(&formatFlag, "format", "", "formattedConfig for stages display, used in conjunction wiht --show-stages")
	rootCmd.AddCommand(command)
}

// common code for create/update/apply and delete actions
func execute(action string, cmd *cobra.Command, args []string) {
	setDefaultProjectID()
	cmd.Printf("%s deployment(s), configs %v, project %s\n", action, args, projectFlag)
	configs := loadConfigs(args)
	stages, err := deployment.Order(configs)
	if err != nil {
		log.Fatalf("Error ordering deployments in dependency order: %v", err)
	}

	isDelete := action == deployment.ActionDelete
	if isDelete {
		// reverse order, dependent goes first for deletion
		for i := len(stages)/2 - 1; i >= 0; i-- {
			opp := len(stages) - 1 - i
			stages[i], stages[opp] = stages[opp], stages[i]
		}
	}
	log.Printf("Ordered dependencies: %v", stages)

	if showStagesFlag {
		showStages(stages)
	} else {
		executeStages(action, stages, isDelete)
	}
}

func showStages(stages [][]deployment.Config) {
	switch formatFlag {
	case "":
		for i, level := range stages {
			log.Printf("---------- Stage %d ----------\n", i)
			for _, config := range level {
				log.Printf("- project: %s, deployment: %s, source: %s", config.GetProject(), config.Name, config.Source())
			}
		}
		log.Printf("------------------------------")
	case "yaml":
		output, _ := yaml.Marshal(formattedConfig(stages))
		log.Println("\n" + string(output))
	case "json":
		output, _ := json.MarshalIndent(formattedConfig(stages), "", "    ")
		log.Println("\n" + string(output))
	}
}

func executeStages(action string, stages [][]deployment.Config, isDelete bool) map[string]map[string]interface{} {
	outputs := make(map[string]map[string]interface{})
	for i, level := range stages {
		for _, config := range level {
			dep := deployment.NewDeployment(config, outputs, !isDelete)
			log.Printf("---------- Stage %d ----------\n", i)
			output, err := dep.Execute(action, previewFlag)
			log.Print(output)
			if err != nil {
				if action == deployment.ActionDelete {
					status, _ := deployment.GetStatus(dep)
					if status == deployment.NotFound {
						// for Delete action, Deployment might not exist, in this case just skip
						log.Printf("Deployment %v does not exist, skip deletion\n", dep)
						continue
					}
				}
				log.Fatalf("Error %s deployment: %v, erro: %v", action, dep, err)
			}
			if validateFlag && (action == deployment.ActionCreate ||
				action == deployment.ActionUpdate ||
				action == deployment.ActionApply) {
				validated, err := validation.ValidateDeployment(config.Name, policyPathFlag, config.GetProject())
				if err != nil {
					log.Fatalf("Error %s validating deployment: %v, erro: %v", action, dep, err)
				}
				if !validated {
					log.Fatalf("Error %s validating deployment: %v", action, dep)
				}
			}

			if previewFlag {
				choice := deployment.GetUserInput("Update(u), Skip (s), or Abort(a) Deployment?", []string{"u", "s", "a"}, os.Stdin)
				switch choice {
				case "u":
					output, err := deployment.ApplyPreview(dep)
					log.Print(output)
					if err != nil {
						log.Fatalf("error %s deployment: %v, error: %v", action, dep, err)
					}
				case "s":
					_, err = deployment.CancelPreview(dep)
					if err != nil {
						log.Fatalf("error cancel-preuvew action for deployment: %v, error: %v", dep, err)
					}
					log.Printf("canceled %s action for deployment: %v", action, dep)
					if action == deployment.ActionCreate {
						log.Printf("delete canceled creation for deployment: %v", dep)
						_, err = deployment.Delete(dep, false)
						if err != nil {
							log.Fatalf("error cancel-preview action for deployment: %v, error: %v", dep, err)
						}
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
	return outputs
}

// listConfigs search for config files according rules described in CLI usage section of following doc:
// https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/dm/docs/userguide.md#syntax
// listConfigs returns map[fileName: fileContent] for files, and list of strings for YAMLs passed as string parameters to CLI
func listConfigs(args []string, errs map[string]error) (map[string]string, []string) {
	resFiles := map[string]string{}
	var resYamls []string
	for _, arg := range args {
		stat, err := os.Stat(arg)
		if os.IsNotExist(err) {
			if deployment.IsYAML(arg) {
				resYamls = append(resYamls, arg)
			} else {
				// check Glob
				matches, err := filepath.Glob(arg)
				if err != nil {
					errs[arg] = fmt.Errorf("error during search files for config: %s, %v", arg, err)
				} else {
					if len(matches) > 0 {
						files, _ := listConfigs(matches, errs)
						resFiles = deployment.AppendMap(resFiles, files)
					} else {
						errs[arg] = fmt.Errorf("no file(s) exist or valid YAML for config param: %s", arg)
					}
				}
			}
		} else if err == nil {
			if stat.IsDir() {
				configs := map[string]string{}
				for _, ext := range supportedExt {
					glob := path.Clean(arg) + "/" + ext
					files, _ := listConfigs([]string{glob}, errs)
					configs = deployment.AppendMap(configs, files)
				}
				resFiles = deployment.AppendMap(resFiles, configs)
				if len(configs) == 0 {
					errs[arg] = fmt.Errorf("no %s files found in directory: %s", strings.Join(supportedExt, ", "), arg)
				}
			} else {
				data, err := ioutil.ReadFile(arg)
				if err != nil {
					log.Fatalf("file %s read error: %v", arg, err)
				}
				resFiles[arg] = string(data)
			}
		} else {
			log.Fatalf("file %s stat error: %v", arg, err)
		}
	}
	return resFiles, resYamls
}

// listConfigs accept list of config parameters (file/directory paths, glob pattern, yaml strings)
// search all possible files in case of directory/glob patterns with listConfigs function and create
// Config objects from loaded data
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

// setDefaultProjectID set deployment.DefaultProjectID variable by search following options:
// The --project command-line option.
// The CLOUD_FOUNDATION_PROJECT_ID environment variable.
// The "default project" configured with the GCP SDK.
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
			log.Print("warning: can't set default project ID from --project arg, CLOUD_FOUNDATION_PROJECT_ID env variable and gcloud default")
		}
		deployment.DefaultProjectID = gcloudDefault
	}
}

// returns anonymous struct suitable for pretty print of JSON and YAML objects
func formattedConfig(stages [][]deployment.Config) interface{} {
	type formattedConfig struct {
		Project    string
		Deployment string
		Source     string
	}
	var configs [][]formattedConfig
	for _, stage := range stages {
		var level []formattedConfig
		for _, config := range stage {
			level = append(level, formattedConfig{config.GetProject(), config.Name, config.Source()})
		}
		configs = append(configs, level)
	}
	return configs
}
