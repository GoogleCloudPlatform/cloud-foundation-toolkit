package main

import (
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/topo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
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

func (c Config) String() string {
	return fmt.Sprintf("%s.%s", c.Project, c.Name)
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
	matches := pattern.FindAllStringSubmatch(c.yaml_string, -1)
	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = match[1]
	}
	if len(result) == 0 {
		return nil
	}
	return result
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
		result := []Config{}
		for dependency := range dependencies {
			result = append(result, dependency)
		}
		return result
	}
	return nil
}

func parseOutRef(text string) (string, string, string, string) {
	array := strings.Split(text, ".")
	return array[0], array[1], array[2], array[3]
}

type Dependency struct {
	from Config
	to   Config
}

func (d Dependency) From() graph.Node {
	return d.from
}

func (d Dependency) To() graph.Node {
	return d.to
}

func (d Dependency) ReversedEdge() graph.Edge {
	return nil
}

type NodesIterator struct {
	nodes []graph.Node
	index int
}

func (i *NodesIterator) Node() graph.Node {
	node := i.nodes[i.index]
	i.index++
	return node
}

func (i *NodesIterator) Len() int {
	return len(i.nodes)
}

func (i *NodesIterator) Reset() {
	i.index = 0
}

func (i *NodesIterator) Next() bool {
	return i.index < i.Len()
}

type ConfigGraph struct {
	nodes     []graph.Node
	edgesFrom map[int64][]Dependency
	edgesTo   map[int64][]Dependency
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
	return &NodesIterator{nodes: c.nodes}
}

func (c ConfigGraph) HasEdgeFromTo(fromid, toid int64) bool {
	nodes := c.From(fromid)
	for nodes.Next() {
		if nodes.Node().ID() == toid {
			return true
		}
	}
	return false
}

func (c ConfigGraph) HasEdgeBetween(xid, yid int64) bool {
	return c.HasEdgeFromTo(xid, yid) || c.HasEdgeFromTo(yid, xid)
}

func (c ConfigGraph) Edge(uid, vid int64) graph.Edge {
	for _, edge := range c.edgesFrom[uid] {
		if edge.To().ID() == vid {
			return edge
		}
	}
	return nil
}

func (c ConfigGraph) sort() (sorted []graph.Node, err error) {
	return topo.Sort(c)
}

func (c *ConfigGraph) AddNode(config Config) {
	c.nodes = append(c.nodes, config)
}

func (c ConfigGraph) From(id int64) graph.Nodes {
	edges := c.edgesFrom[id]
	res := []graph.Node{}
	if edges != nil {
		for _, edge := range edges {
			res = append(res, edge.To())
		}
	}
	return &NodesIterator{nodes: res}
}

func (c ConfigGraph) To(id int64) graph.Nodes {
	edges := c.edgesTo[id]
	res := []graph.Node{}
	if edges != nil {
		for _, edge := range edges {
			res = append(res, edge.From())
		}
	}
	return &NodesIterator{nodes: res}
}

func (c *ConfigGraph) AddDependency(from Config, to Config) {
	dependencies := c.edgesFrom[from.ID()]
	if dependencies == nil {
		dependencies = []Dependency{}
	}
	dep := &Dependency{from: from, to: to}
	c.edgesFrom[from.ID()] = append(dependencies, *dep)
	c.edgesTo[to.ID()] = append(dependencies, *dep)
}

func NewConfigGraph(configs []Config) *ConfigGraph {
	instance := &ConfigGraph{}
	instance.edgesFrom = make(map[int64][]Dependency)
	instance.edgesTo = make(map[int64][]Dependency)

	for _, config := range configs {
		instance.AddNode(config)
	}

	for _, config := range configs {
		deps := config.findAllDependencies(configs)
		for _, dep := range deps {
			instance.AddDependency(config, dep)
		}
	}
	return instance
}
