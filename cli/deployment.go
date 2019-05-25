package main

import (
	"gonum.org/v1/gonum/graph"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
)

var /* const */ pattern = regexp.MustCompile(`\$\(out\.(?P<token>[-.a-zA-Z0-9]+)\)`)

type Config struct {
	Name        string
	Project     string
	file_path   string
	yaml_string string
}

// implementation of graph.Node interface
func (c Config) ID() int64 {
	return hash64(c.Project + "." + c.Name)
}

func NewConfig(file_path string) *Config {
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		log.Fatal(err)
	}
	yaml_string := string(data)

	config := &Config{
		file_path:   file_path,
		yaml_string: yaml_string,
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config
}

func (c Config) findAllOutRefs() []string {
	return pattern.FindAllString(c.yaml_string, -1)
}

type Dependency struct {
	from Config
	to   Config
}

func (d Dependency) From() Config {
	return d.from
}

func (d Dependency) To() Config {
	return d.to
}

type ConfigGraphIterator struct {
	graph *ConfigGraph
	index int
}

func (i *ConfigGraphIterator) Node() graph.Node {
	return i.graph.nodes[i.index]
}

func (i *ConfigGraphIterator) Len() int {
	return len(i.graph.nodes)
}

func (i *ConfigGraphIterator) Reset() {
	i.index = 0
}

func (i *ConfigGraphIterator) Next() bool {
	return i.index < i.Len()-1
}

type ConfigGraph struct {
	nodes     []Config
	edgesFrom map[int64]Dependency
	edgesTo   map[int64]Dependency
}

func (c ConfigGraph) Node(id int64) graph.Node {
	for _, node := range c.nodes {
		if node.ID() == id {
			return node
		}
	}
	return nil
}

func (c ConfigGraph) Nodes() graph.Nodes {
	return &ConfigGraphIterator{graph: &c}
}

func (c ConfigGraph) From(id int64) graph.Nodes {
	return c.edgesFrom[id]
}

func NewConfigGraph(config []Config) {
	graph := &ConfigGraph{}

}
