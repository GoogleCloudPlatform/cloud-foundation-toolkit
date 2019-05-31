package deployment

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
	"strings"
)

var /* const */ pattern = regexp.MustCompile(`\$\(out\.(?P<token>[-.a-zA-Z0-9]+)\)`)

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

func NewConfig(data string, file string) *Config {
	config := &Config{
		file: file,
		data: data,
	}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config
}

func (c Config) findAllDependencies(configs []Config) []Config {
	refs := c.findAllOutRefs()
	if refs != nil {
		dependencies := map[int64]Config{}
		for _, ref := range refs {
			var dependency *Config
			project, deployment, _, _ := parseOutRef(ref)
			for i, config := range configs {
				if config.Project == project && config.Name == deployment {
					dependency = &configs[i]
					fmt.Printf("Found dependency %s\n", dependency.String())
				}
			}
			if dependency == nil {
				log.Fatalf("Could not find config for project = %s, deployment = %s", project, deployment)
			}
			dependencies[dependency.ID()] = *dependency
		}
		var result []Config
		for _, dependency := range dependencies {
			fmt.Printf("Adding dependency %s\n", dependency.String())
			result = append(result, dependency)
		}
		return result
	}
	return nil
}

func (c Config) Yaml() ([]byte, error) {
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

func (c Config) importsAbsolutePath() (imports interface{}, typeMap map[string]string) {
	typeMap = map[string]string{}
	for i, imp := range c.Imports {
		newPath := AbsolutePath(c.file, imp.Path)
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

// implementation of graph.Node interface
func (c Config) ID() int64 {
	return hash64(c.Project + "." + c.Name)
}

func (c Config) String() string {
	return fmt.Sprintf("%s.%s", c.Project, c.Name)
}

func (c Config) findAllOutRefs() []string {
	matches := pattern.FindAllStringSubmatch(c.data, -1)
	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = match[1]
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func parseOutRef(text string) (project string, deployment string,
	resource string, property string) {
	array := strings.Split(text, ".")
	return array[0], array[1], array[2], array[3]
}
