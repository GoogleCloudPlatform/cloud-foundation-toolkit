package deployment

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// Pattern to parse $(project.deployment.resource.name) and $(deployment.resource.name).
var pattern = regexp.MustCompile(`\$\(out\.(?P<token>[-.a-zA-Z0-9]+)\)`)

// Config struct keep config data parsed from passed config YAML.
type Config struct {
	// Name field contains Deployment Name it could be initialized from
	// config YAML field and if it missing in YAML, from config YAML file name itself.
	Name string
	// Project field contains GCP Project Id it could be initialized from config YAML, env variable, gcloud default.
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
	data      string
}

// The Config type represents configuration data parsed from YAML.
func NewConfig(data string, file string) Config {
	config := Config{
		file: file,
		data: data,
	}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config
}

// FullName returns a name for the configuration that is intended to be unique.
// The name is in the format "ProjectName.DeploymentName" and may be used as a map key.
func (c Config) FullName() string {
	return c.Project + "." + c.Name
}

// String implements fmt.Stringer for the Config type.
func (c Config) String() string {
	return c.FullName()
}

// The YAML function marshals a Config object to YAML. It sets all relative import paths to absolute paths and removes
// all custom elements that the gcloud deployment manager is not aware of, including name, project, and description.
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

// Function findAllDependencies finds all dependencies base on cross-deployment references found in config YAML,
// if no dependencies found, returns empty slice.
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
		fmt.Printf("Adding dependency %s -> %s\n", c.FullName(), dependency.FullName())
		result = append(result, dependency)
	}
	return result
}

func (c Config) importsAbsolutePath() (imports interface{}, typeMap map[string]string) {
	typeMap = map[string]string{}
	for i, imp := range c.Imports {
		newPath := ReparentPath(c.file, imp.Path)
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
		outputsMap, err := GetOutputs(arr[0], arr[1])
		if err != nil {
			log.Fatalf("Erorr getting outputs for deployment: %s, error: %v", fullName, err)
		}
		outputs[fullName] = outputsMap
	}
	value, ok := outputsMap[res+"."+name]
	fullRef := fmt.Sprintf("$(out.%s)", ref)
	if !ok {
		log.Fatalf("Could not resolve reference: %s", fullRef)
	}
	return value
}

func parseOutRef(text string) (fullName string, resource string, property string) {
	array := strings.Split(text, ".")
	return array[0] + "." + array[1], array[2], array[3]
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
