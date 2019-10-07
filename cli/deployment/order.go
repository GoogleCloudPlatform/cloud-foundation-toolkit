package deployment

import (
	"errors"
	"fmt"
	"log"
)

// Order function receives map of configs with Config.FullName() string as a key,
// find dependencies between them, and order them topologically using directed graph.
// Returns array of arrays of config, each inner array represent configs that could be created in parallel
// each next level depends on previous
func Order(configs map[string]Config) ([][]Config, error) {
	var nodes []string
	// we don't know number or dependencies, so initial size is 0
	edges := make([]edge, 0)
	for _, config := range configs {
		nodes = append(nodes, config.FullName())
		deps, err := config.findAllDependencies(configs)
		if err != nil {
			log.Fatalf("Failed to find dependencies for config = %v", config)
		}
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
	var res [][]Config
	for _, level := range sorted {
		var resLevel []Config
		for _, name := range level {
			resLevel = append(resLevel, configs[name])
		}
		res = append(res, resLevel)
	}
	return res, nil
}

// directedGraph is a graph representation of Depoyment configs, each node in graph represent Deployment config object,
// each edge of the graph reflects cross referenced dependency, for example A -> B, means B depends on A and A should be created first.
type directedGraph struct {
	nodes         []string
	outgoingNodes map[string]map[string]int
	incomingNodes map[string]int
}

// Function newDirectedGraph create and initialize new directedGraph instance.
// If nodes parameter will contain duplicate values, nil and error will be returned.
// If edges parameter will contain non existing from node, nil and erro will be returned.
func newDirectedGraph(nodes []string, edges []edge) (*directedGraph, error) {
	g := &directedGraph{
		nodes:         make([]string, 0, len(nodes)),
		incomingNodes: make(map[string]int),
		outgoingNodes: make(map[string]map[string]int),
	}

	for _, node := range nodes {
		g.nodes = append(g.nodes, node)
		if _, ok := g.outgoingNodes[node]; ok {
			return nil, fmt.Errorf("node %s already added to graph", node)
		}
		g.outgoingNodes[node] = make(map[string]int)
		g.incomingNodes[node] = 0
	}

	for _, edge := range edges {
		node, ok := g.outgoingNodes[edge.from]
		if !ok {
			return nil, fmt.Errorf("no node %s exists in graph", edge.from)
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

// Function topologicalSort() contains core logic of topological search here is to find root nodes (no incoming nodes),
// remove them from graph and repeat until all grahp will be traversed.
func (g *directedGraph) topologicalSort() ([][]string, error) {
	var result [][]string

	rootNodes := make([]string, 0, len(g.nodes))
	for _, n := range g.nodes {
		if g.incomingNodes[n] == 0 {
			rootNodes = append(rootNodes, n)
		}
	}
	result = append(result, rootNodes)

	var nextLevel []string
	for len(rootNodes) > 0 {
		var current string
		current, rootNodes = rootNodes[0], rootNodes[1:]

		outgoingNodes := make([]string, len(g.outgoingNodes[current]))
		for outgoingNode, i := range g.outgoingNodes[current] {
			outgoingNodes[i-1] = outgoingNode
		}

		for _, outgoingNode := range outgoingNodes {
			g.unsafeRemoveEdge(current, outgoingNode)

			if g.incomingNodes[outgoingNode] == 0 {
				nextLevel = append(nextLevel, outgoingNode)
			}
		}

		if len(rootNodes) == 0 && len(nextLevel) > 0 {
			for _, next := range nextLevel {
				rootNodes = append(rootNodes, next)
			}
			result = append(result, nextLevel)
			nextLevel = nil
		}
	}

	outgoingCount := 0
	for _, v := range g.incomingNodes {
		outgoingCount += v
	}

	if outgoingCount > 0 {
		return nil, errors.New("cycle detected in graph")
	}
	return result, nil
}
