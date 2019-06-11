package deployment

import (
	"errors"
	"fmt"
	"log"
)

// Order function receives map of configs with Config.FullName() string as a key,
// find dependencies between them, and order them topologically using directed graph.
// Returns array of ordered config's objects and error if configs have circle dependencies,
func Order(configs map[string]Config) ([]Config, error) {
	size := len(configs)

	nodes := make([]string, 0)
	// we don't know number or dependencies, so initial size is 0
	edges := make([]edge, 0)
	for _, config := range configs {
		nodes = append(nodes, config.FullName())
		deps := config.findAllDependencies(configs)
		for _, dep := range deps {
			edges = append(edges, edge{
				from: dep.FullName(),
				to:   config.FullName(),
			})
		}
	}

	graph, err := newDirectedGraph(nodes, edges)
	if err != nil {
		log.Printf("error creating dependecy grahp: %v", err)
		return nil, err
	}

	sorted, err := graph.topologicalSort()
	if err != nil {
		log.Printf("error ordering configs: %v", err)
		return nil, err
	}
	res := make([]Config, size)
	for i, name := range sorted {
		res[i] = configs[name]
	}
	return res, nil
}

// directedGraph struct used to build graph of cross-config depencencies
// and do topological sort to define ordering of deployment creation
type directedGraph struct {
	nodes         []string
	outgoingNodes map[string]map[string]int
	incomingNodes map[string]int
}

// newDirectedGraph create and initialize new directedGraph instance
// if nodes parameter will contain duplicate values, nil and error will be returned
// if edges parameter will contain non existing from node, nil and erro will be returned
func newDirectedGraph(nodes []string, edges []edge) (*directedGraph, error) {
	g := &directedGraph{
		nodes:         make([]string, 0, len(nodes)),
		incomingNodes: make(map[string]int),
		outgoingNodes: make(map[string]map[string]int),
	}

	for _, node := range nodes {
		g.nodes = append(g.nodes, node)
		if _, ok := g.outgoingNodes[node]; ok {
			return nil, errors.New(fmt.Sprintf("node %s already added to graph", node))
		}
		g.outgoingNodes[node] = make(map[string]int)
		g.incomingNodes[node] = 0
	}

	for _, edge := range edges {
		node, ok := g.outgoingNodes[edge.from]
		if !ok {
			return nil, errors.New(fmt.Sprintf("no node %s exists in graph", edge.from))
		}

		node[edge.to] = len(node) + 1
		g.incomingNodes[edge.to]++
	}
	return g, nil
}

type edge struct {
	from string
	to   string
}

func (g *directedGraph) unsafeRemoveEdge(from, to string) {
	delete(g.outgoingNodes[from], to)
	g.incomingNodes[to]--
}

// main logic of topological search here is to find root nodes (no incoming nodes),
// remove them from graph and repeat untill all grahp will be traversed
func (g *directedGraph) topologicalSort() ([]string, error) {
	sorted := make([]string, 0, len(g.nodes))
	rootNodes := make([]string, 0, len(g.nodes))

	for _, n := range g.nodes {
		if g.incomingNodes[n] == 0 {
			rootNodes = append(rootNodes, n)
		}
	}

	for len(rootNodes) > 0 {
		var current string
		current, rootNodes = rootNodes[0], rootNodes[1:]
		sorted = append(sorted, current)

		outgoingNodes := make([]string, len(g.outgoingNodes[current]))
		for outgoingNode, i := range g.outgoingNodes[current] {
			outgoingNodes[i-1] = outgoingNode
		}

		for _, outgoingNode := range outgoingNodes {
			g.unsafeRemoveEdge(current, outgoingNode)

			if g.incomingNodes[outgoingNode] == 0 {
				rootNodes = append(rootNodes, outgoingNode)
			}
		}
	}

	outgoingCount := 0
	for _, v := range g.incomingNodes {
		outgoingCount += v
	}

	if outgoingCount > 0 {
		return nil, errors.New("cycle detected in graph")
	}

	return sorted, nil
}
