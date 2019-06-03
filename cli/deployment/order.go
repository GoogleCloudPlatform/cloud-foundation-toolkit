package deployment

import (
	"errors"
	"log"
)

func Order(configs map[string]Config) ([]Config, error) {
	size := len(configs)
	graph := NewDirectedGraph(size)

	for _, config := range configs {
		graph.AddNode(config.FullName())
	}

	for _, config := range configs {
		deps := config.findAllDependencies(configs)
		for _, dep := range deps {
			graph.AddEdge(dep.FullName(), config.FullName())
		}
	}

	sorted, err := graph.sortTopologically()
	if err != nil {
		log.Printf("Error duing ordering configs dependencies %v", err)
		return nil, err
	}
	res := make([]Config, size)
	for i, name := range sorted {
		res[i] = configs[name]
	}
	return res, nil
}

type DirectedGraph struct {
	nodes          []string
	outgoingNodes  map[string]map[string]int
	incommingNodes map[string]int
}

func NewDirectedGraph(cap int) *DirectedGraph {
	return &DirectedGraph{
		nodes:          make([]string, 0, cap),
		incommingNodes: make(map[string]int),
		outgoingNodes:  make(map[string]map[string]int),
	}
}

func (g *DirectedGraph) AddNode(name string) bool {
	g.nodes = append(g.nodes, name)

	if _, ok := g.outgoingNodes[name]; ok {
		return false
	}
	g.outgoingNodes[name] = make(map[string]int)
	g.incommingNodes[name] = 0
	return true
}

func (g *DirectedGraph) AddNodes(names ...string) bool {
	for _, name := range names {
		if ok := g.AddNode(name); !ok {
			return false
		}
	}
	return true
}

func (g *DirectedGraph) AddEdge(from, to string) bool {
	node, ok := g.outgoingNodes[from]
	if !ok {
		return false
	}

	node[to] = len(node) + 1
	g.incommingNodes[to]++

	return true
}

func (g *DirectedGraph) RemoveEdge(from, to string) bool {
	// check if edge exists
	if _, ok := g.outgoingNodes[from]; !ok {
		return false
	}
	g.unsafeRemoveEdge(from, to)
	return true
}

func (g *DirectedGraph) unsafeRemoveEdge(from, to string) {
	delete(g.outgoingNodes[from], to)
	g.incommingNodes[to]--
}

func (graph *DirectedGraph) sortTopologically() ([]string, error) {
	sorted := make([]string, 0, len(graph.nodes))
	rootNodes := make([]string, 0, len(graph.nodes))

	for _, n := range graph.nodes {
		if graph.incommingNodes[n] == 0 {
			rootNodes = append(rootNodes, n)
		}
	}

	for len(rootNodes) > 0 {
		var current string
		current, rootNodes = rootNodes[0], rootNodes[1:]
		sorted = append(sorted, current)

		outgoingNodes := make([]string, len(graph.outgoingNodes[current]))
		for outgoingNode, i := range graph.outgoingNodes[current] {
			outgoingNodes[i-1] = outgoingNode
		}

		for _, outgoingNode := range outgoingNodes {
			graph.unsafeRemoveEdge(current, outgoingNode)

			if graph.incommingNodes[outgoingNode] == 0 {
				rootNodes = append(rootNodes, outgoingNode)
			}
		}
	}

	outgoingCount := 0
	for _, v := range graph.incommingNodes {
		outgoingCount += v
	}

	if outgoingCount > 0 {
		return nil, errors.New("Circle dependencies in graph")
	}

	return sorted, nil
}
