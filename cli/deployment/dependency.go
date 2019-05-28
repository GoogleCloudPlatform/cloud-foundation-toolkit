package deployment

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/topo"
)

type Dependency struct {
	from Config
	to   Config
}

type DependencyGraph struct {
	nodes     []graph.Node
	edgesFrom map[int64][]Dependency
	edgesTo   map[int64][]Dependency
}

func NewDependencyGraph(configs []Config) *DependencyGraph {
	instance := &DependencyGraph{}
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

func (c DependencyGraph) Node(id int64) graph.Node {
	for _, node := range c.nodes {
		if node.ID() == id {
			return node
		}
	}
	return nil
}

func (c DependencyGraph) Nodes() graph.Nodes {
	return &NodesIterator{nodes: c.nodes}
}

func (c DependencyGraph) HasEdgeFromTo(fromid, toid int64) bool {
	nodes := c.From(fromid)
	for nodes.Next() {
		if nodes.Node().ID() == toid {
			return true
		}
	}
	return false
}

func (c DependencyGraph) HasEdgeBetween(xid, yid int64) bool {
	return c.HasEdgeFromTo(xid, yid) || c.HasEdgeFromTo(yid, xid)
}

func (c DependencyGraph) Edge(uid, vid int64) graph.Edge {
	for _, edge := range c.edgesFrom[uid] {
		if edge.To().ID() == vid {
			return edge
		}
	}
	return nil
}

func (c DependencyGraph) Order() (ordered []Config, err error) {
	nodes, err := topo.Sort(c)
	if err != nil {
		return nil, err
	}
	res := make([]Config, len(nodes))
	for i, node := range nodes {
		res[i] = node.(Config)
	}
	return res, err
}

func (c *DependencyGraph) AddNode(config Config) {
	c.nodes = append(c.nodes, config)
}

func (c DependencyGraph) From(id int64) graph.Nodes {
	edges := c.edgesFrom[id]
	var res []graph.Node
	if edges != nil {
		for _, edge := range edges {
			res = append(res, edge.To())
		}
	}
	return &NodesIterator{nodes: res}
}

func (c DependencyGraph) To(id int64) graph.Nodes {
	edges := c.edgesTo[id]
	var res []graph.Node
	if edges != nil {
		for _, edge := range edges {
			res = append(res, edge.From())
		}
	}
	return &NodesIterator{nodes: res}
}

func (c *DependencyGraph) AddDependency(from Config, to Config) {
	dependencies := c.edgesFrom[from.ID()]
	if dependencies == nil {
		dependencies = []Dependency{}
	}
	dep := &Dependency{from: from, to: to}
	c.edgesFrom[from.ID()] = append(dependencies, *dep)
	c.edgesTo[to.ID()] = append(dependencies, *dep)
}
