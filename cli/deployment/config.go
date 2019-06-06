package deployment

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// patern to parse $(project.deployment.resource.name) and $(deployment.resource.name)
var pattern = regexp.MustCompile(`\$\(out\.(?P<token>[-.a-zA-Z0-9]+)\)`)

// Config struct keep config data parsed from passed config.yaml
type Config struct {
	Name        string
	Project     string
	Description string
	Imports     []struct {
		Name string
		Path string
	}
	Resources []interface{}
	file      string
	data      string
}

// NewConfig creates new Config object from provided yaml file
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

// FullName returns name in form of ProjectName.DeploymentName, this name should be unique and it could be used as map key
// for maps like map[string]Config
func (c Config) FullName() string {
	return c.Project + "." + c.Name
}

func (c Config) String() string {
	return c.FullName()
}

// YAML function converts Config object to yaml
// overrides all relative paths for imports to absolute form,
// removes all custom elements gcloud deployment manager not aware of (name, project, description)
func (c Config) YAML() ([]byte, error) {
	imports, typeMap := c.importsAbsolutePath()

	tmp := struct {
		Imports   interface{}
		Resources interface{}
	}{
		Imports:   imports,
		Resources: c.resources(typeMap),
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

func (c Config) resources(typeMap map[string]string) []interface{} {
	if len(typeMap) > 0 {
		for i := range c.Resources {
			res := c.Resources[i].(map[interface{}]interface{})
			if typeMap[res["type"].(string)] != "" {
				res["type"] = typeMap[res["type"].(string)]
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

func parseOutRef(text string) (fullName string, resource string, property string) {
	array := strings.Split(text, ".")
	return array[0] + "." + array[1], array[2], array[3]
}
