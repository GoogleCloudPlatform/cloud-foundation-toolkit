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
	Name    string
	Project string
	file    string
	data    string
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
		dependencies := map[Config]bool{}
		for _, ref := range refs {
			var dependency *Config
			project, deployment, _, _ := parseOutRef(ref)
			for _, config := range configs {
				if config.Project == project && config.Name == deployment {
					dependency = &config
				}
			}
			if dependency == nil {
				log.Fatalf("Could not find config for project = %s, deployment = %s", project, deployment)
			}
			dependencies[*dependency] = true
		}
		var result []Config
		for dependency := range dependencies {
			result = append(result, dependency)
		}
		return result
	}
	return nil
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
	resouce string, property string) {
	array := strings.Split(text, ".")
	return array[0], array[1], array[2], array[3]
}
