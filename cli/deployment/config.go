package deployment

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var DefaultProjectID string

// patern to parse $(project.deployment.resource.name) and $(deployment.resource.name)
var pattern = regexp.MustCompile(`\$\(out\.(?P<token>[-.a-zA-Z0-9]+)\)`)

// Config struct keep config data parsed from passed config.yaml
type Config struct {
	Name string
	// don't use this varible directly, use GetProjectID instead!!!
	Project     string
	Description string
	Imports     []struct {
		Name string
		Path string
	}
	Resources []map[string]interface{}
	file      string
	dir       string
	data      string
}

// NewConfig creates new Config object from provided yaml file
func NewConfig(data string, file string) Config {
	config := Config{
		file: file,
		data: data,
	}

	if len(file) > 0 {
		config.file = filepath.Clean(file)
		config.dir = filepath.Dir(file)
	} else {
		// yaml passed as sting through parameter
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
	// check project id and deployment name
	if len(config.Project) == 0 {
		config.Project = DefaultProjectID
	}
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

// returns project id if set in config file, otherwise default
func (c Config) GetProject() string {
	if len(c.Project) > 0 {
		return c.Project
	} else {
		return DefaultProjectID
	}
}

// returns file path in case of file source or yaml string in case of yaml source
func (c Config) Source() string {
	if len(c.file) > 0 {
		return c.file
	} else {
		return c.data
	}
}

// FullName returns name in form of ProjectName.DeploymentName, this name should be unique and it could be used as map key
// for maps like map[string]Config
func (c Config) FullName() string {
	return c.GetProject() + "." + c.Name
}

func (c Config) String() string {
	return c.FullName()
}

// YAML function converts Config object to yaml
// overrides all relative paths for imports to absolute form,
// removes all custom elements gcloud deployment manager not aware of (name, project, description)
func (c Config) YAML(outputs map[string]map[string]interface{}) ([]byte, error) {
	imports, typeMap := c.importsAbsolutePath()
	resources := c.resources(typeMap)
	for _, value := range resources {
		replaceOutRefsResource(value, outputs)
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

// findAllDependencies finds all dependencies base on cross-deployment references found in config yaml,
// if no dependencies found, returns empty slice
func (c Config) findAllDependencies(configs map[string]Config) []Config {
	var result []Config
	refs := c.findAllOutRefs()
	dependencies := map[string]Config{}
	for _, ref := range refs {
		fullName, _, _ := parseOutRef(ref)
		dependency, found := configs[fullName]
		if !found {
			log.Fatalf("Could not find config for deployment = %s", fullName)
		}
		dependencies[fullName] = dependency
	}

	for _, dependency := range dependencies {
		log.Printf("Adding dependency %s -> %s\n", c.FullName(), dependency.FullName())
		result = append(result, dependency)
	}
	return result
}

func (c Config) importsAbsolutePath() (imports interface{}, typeMap map[string]string) {
	typeMap = map[string]string{}
	for i, imp := range c.Imports {
		newPath := AbsolutePath(c.dir, imp.Path)
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

func getOutRefValue(ref string, outputs map[string]map[string]interface{}) interface{} {
	fullName, res, name := parseOutRef(ref)
	outputsMap := outputs[fullName]
	if outputsMap == nil {
		arr := strings.Split(fullName, ".")
		var err error
		outputsMap, err = GetOutputs(arr[0], arr[1])
		if err != nil {
			log.Fatalf("Erorr getting outputs for deployment: %s, error: %v", fullName, err)
		}
		outputs[fullName] = outputsMap
	}
	value, ok := outputsMap[res+"."+name]
	fullRef := fmt.Sprintf("$(out.%s)", ref)
	if !ok {
		log.Fatalf("Unresolved dependency: %s. Deployment: %s , on whichother resources depended, neither was specifiedin the submitted congigs nor existed inDeployment Manager", fullRef, fullName)
	}
	return value
}

func parseOutRef(text string) (fullName string, resource string, property string) {
	array := strings.Split(text, ".")
	project, deploymentName, resource, property := DefaultProjectID, array[0], array[1], array[2]
	if len(array) == 4 {
		project, deploymentName, resource, property = array[0], array[1], array[2], array[3]
	}
	return project + "." + deploymentName, resource, property
}

func replaceOutRefsResource(resource interface{}, outputs map[string]map[string]interface{}) interface{} {
	switch t := resource.(type) {
	case string:
		value := resource.(string)
		match := pattern.FindStringSubmatch(value)
		if match == nil {
			return value
		} else {
			return getOutRefValue(match[1], outputs)
		}
	case map[string]interface{}:
		values := resource.(map[string]interface{})
		for key, value := range values {
			values[key] = replaceOutRefsResource(value, outputs)
		}
		return values
	case map[interface{}]interface{}:
		values := resource.(map[interface{}]interface{})
		for key, value := range values {
			values[key] = replaceOutRefsResource(value, outputs)
		}
		return values
	case []interface{}:
		values := resource.([]interface{})
		var result = []interface{}{}
		for _, value := range values {
			result = append(result, replaceOutRefsResource(value, outputs))
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
