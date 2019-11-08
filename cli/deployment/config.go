package deployment

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// DefaultProjectID holds default GCP project ID, it used in case project ID not provided in configuration file,
// see common.setDefaultProjectID for variable initialization details.
var DefaultProjectID string

// Pattern to parse $(project.deployment.resource.name) and $(deployment.resource.name).
var pattern = regexp.MustCompile(`\$\(out\.(?P<token>[-.a-zA-Z0-9]+)\)`)

// Config struct keep config data parsed from passed config YAML.
type Config struct {
	// Name field contains Deployment Name it could be initialized from
	// config YAML field and if it missing in YAML, from config YAML file name itself.
	Name string
	// Project field contains GCP Project Id it could be initialized from config YAML, env variable, gcloud default.
	// Note: don't access this member directly, use GetProjectID instead.
	Project string
	// Deployment string contains GCP Deployment description it not required and can be empty.
	Description string
	// Imports struct contains all "import" entries from config YAML file.
	Imports []struct {
		// Name of the import statement.
		Name string
		// Path to import file, can be relative or absolute.
		Path string
	}
	// Resources map contains all Deployment Resources, listed in deployment config YAML file.
	// TODO explain what value used as map key
	Resources []map[string]interface{}
	file      string
	dir       string
	data      string
}

// The Config type represents configuration data parsed from YAML.
func NewConfig(data string, file string) Config {
	config := Config{
		file: file,
		data: data,
	}

	if len(file) > 0 {
		config.file = filepath.Clean(file)
		config.dir = filepath.Dir(file)
	} else {
		// YAML passed as string through parameter
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not get current folder path: %v", err)
		}
		config.dir = dir
	}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("failed unmarshal config yaml %s, error: %v", file, err)
	}
	// check deployment name
	if len(config.Name) == 0 {
		if len(config.file) == 0 {
			// if file empty, than yaml string as a parameter was passed, it SHOULD contain name field
			log.Fatalf("name field not defined in yaml: %s", data)
		} else {
			config.Name = DeploymentNameFromFile(config.file)
		}
	}
	return config
}

// GetProject returns project id if it set in config file, otherwise default value defined in DefaultProjectID var.
func (c Config) GetProject() string {
	if len(c.Project) > 0 {
		return c.Project
	} else if len(DefaultProjectID) > 0 {
		return DefaultProjectID
	} else {
		log.Fatal("warning: can't set default project ID from --project arg, CLOUD_FOUNDATION_PROJECT_ID env variable and gcloud default")
		return ""
	}
}

// Source returns the the path of the file for a file-backed configuration, or a YAML string for YAML configurations.
func (c Config) Source() string {
	if len(c.file) > 0 {
		return c.file
	} else {
		return c.data
	}
}

// FullName returns a name for the configuration that is intended to be unique.
// The name is in the format "ProjectName.DeploymentName" and may be used as a map key.
func (c Config) FullName() string {
	return c.GetProject() + "." + c.Name
}

// String implements fmt.Stringer for the Config type.
func (c Config) String() string {
	return c.FullName()
}

// YAML function marshals a Config object to YAML. It sets all relative import paths to absolute paths and removes
// all custom elements that the gcloud deployment manager is not aware of, including name, project, and description.
// output param contains map of map with Outputs variables of all Deployments created/updated by current cli run,
// where key of first map is project.deployment names and map[string]interface{} - map of Deployment output properties names->values.
// YAML returns byte array of config with output references substituted by real values from outputs parameter map.
func (c Config) YAML(outputs map[string]map[string]interface{}) ([]byte, error) {
	imports, typeMap := c.importsAbsolutePath()
	resources := c.resources(typeMap)
	for _, value := range resources {
		c.replaceOutRefsResource(value, outputs)
	}

	tmp := struct {
		Imports   interface{}
		Resources interface{}
	}{
		Imports:   imports,
		Resources: resources,
	}
	return yaml.Marshal(tmp)
}

// Function findAllDependencies finds all dependencies base on cross-deployment references found in config YAML,
// if no dependencies found, returns empty slice.
func (c Config) findAllDependencies(configs map[string]Config) ([]Config, error) {
	var result []Config
	refs := c.findAllOutRefs()
	dependencies := map[string]Config{}
	for _, ref := range refs {
		fullName, _, _ := c.parseOutRef(ref)
		dependency, found := configs[fullName]
		exists := false
		if !found {
			projectAndName := strings.Split(fullName, ".")
			deployment := &Deployment{
				config: Config{
					Project: projectAndName[0],
					Name:    projectAndName[1],
				},
			}
			status, err := GetStatus(deployment)

			if err != nil {
				log.Printf("GetStatus for deployment = %s", fullName)
				return nil, err
			}
			switch status {
			case Done:
				exists = true
			case NotFound:
				message := fmt.Sprintf("Could not find config or existing deployment = %s", fullName)
				log.Print(message)
				return nil, errors.New(message)
			case Pending, Error, Running:
				message := fmt.Sprintf("Dependency deployment = %sis in %v state", fullName, status)
				log.Print(message)
				return nil, errors.New(message)
			}
		}
		if !exists {
			// if dependency already exists no need to count during creation ordering
			dependencies[fullName] = dependency
		}
	}

	for _, dependency := range dependencies {
		log.Printf("Adding dependency %s -> %s\n", c.FullName(), dependency.FullName())
		result = append(result, dependency)
	}
	return result, nil
}

func (c Config) importsAbsolutePath() (imports interface{}, typeMap map[string]string) {
	typeMap = map[string]string{}
	for i, imp := range c.Imports {
		newPath := ReparentPath(c.dir, imp.Path)
		if newPath != c.Imports[i].Path {
			typeMap[c.Imports[i].Path] = newPath
		}
		c.Imports[i].Path = newPath
	}
	return c.Imports, typeMap
}

func (c Config) resources(typeMap map[string]string) []map[string]interface{} {
	if len(typeMap) > 0 {
		for _, element := range c.Resources {
			if typeMap[element["type"].(string)] != "" {
				element["type"] = typeMap[element["type"].(string)]
			}
		}
	}
	return c.Resources
}

func (c Config) findAllOutRefs() []string {
	var result []string
	matches := pattern.FindAllStringSubmatch(c.data, -1)
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func (c Config) getOutRefValue(ref string, outputs map[string]map[string]interface{}) interface{} {
	fullName, res, name := c.parseOutRef(ref)
	outputsMap := outputs[fullName]
	if outputsMap == nil {
		arr := strings.Split(fullName, ".")
		var err error
		outputsMap, err = GetOutputs(arr[0], arr[1])
		if err != nil {
			log.Fatalf("Error getting outputs for deployment: %s, error: %v", fullName, err)
		}
		outputs[fullName] = outputsMap
	}
	value, ok := outputsMap[res+"."+name]
	if !ok {
		fullRef := fmt.Sprintf("$(out.%s)", ref)
		log.Fatalf("Unresolved dependency: %s. Deployment: %s , on which other resources depended, was neither specified in the submitted configs nor existed in Deployment Manager", fullRef, fullName)
	}
	return value
}

func (c Config) parseOutRef(text string) (fullName string, resource string, property string) {
	array := strings.Split(text, ".")
	project, deploymentName, resource, property := c.GetProject(), array[0], array[1], array[2]
	if len(array) == 4 {
		project, deploymentName, resource, property = array[0], array[1], array[2], array[3]
	}
	return project + "." + deploymentName, resource, property
}

func (c Config) replaceOutRefsResource(resource interface{}, outputs map[string]map[string]interface{}) interface{} {
	switch t := resource.(type) {
	case string:
		value := resource.(string)
		match := pattern.FindStringSubmatch(value)
		if match == nil {
			return value
		} else {
			result := c.getOutRefValue(match[1], outputs)
			if reflect.TypeOf(result).Kind() == reflect.String {
				value = strings.Replace(value, match[0], result.(string), 1);
				return c.replaceOutRefsResource(value, outputs)
			} else {
				return result
			}
		}
	case map[string]interface{}:
		values := resource.(map[string]interface{})
		for key, value := range values {
			values[key] = c.replaceOutRefsResource(value, outputs)
		}
		return values
	case map[interface{}]interface{}:
		values := resource.(map[interface{}]interface{})
		for key, value := range values {
			values[key] = c.replaceOutRefsResource(value, outputs)
		}
		return values
	case []interface{}:
		values := resource.([]interface{})
		var result = []interface{}{}
		for _, value := range values {
			result = append(result, c.replaceOutRefsResource(value, outputs))
		}
		return result
	case bool:
		return resource
	case int:
		return resource
	default:
		log.Fatalf("unexpected yaml element type: %v", t)
	}
	return nil
}
