// Package launchpad file runtime.go contains runtime related support for
// evaluation hierarchy and evaluated object tracking.
package launchpad

import (
	"errors"
	"fmt"
)

// globalState keeps track of evaluation order while parsing YAML and stores
// metadata like configurations.
type globalState struct {
	stack           []stackFrame
	outputDirectory string
	outputFlavor    outputFlavor
	evaluated       struct {
		folders folders
	}
}

// ==== Evaluation Stack ====

// stackable interface determines what can be pushed onto the stack.
//
// stackable is also synonymous to yaml.Unmarshaler.
type stackable interface {
	UnmarshalYAML(unmarshal func(interface{}) error) error
}

// stackFrame defines a single evaluation hierarchy.
type stackFrame struct {
	stackType crdKind
	stackPtr  stackable
}

// push pushes a new stackFrame onto current stack.
func (g *globalState) push(stackType crdKind, stackPtr stackable) {
	g.stack = append(g.stack, stackFrame{
		stackType: stackType,
		stackPtr:  stackPtr,
	})
}

// pop ejects top of stack's stackFrame.
func (g *globalState) pop() (stackFrame, error) {
	l := len(g.stack)
	if l == 0 {
		return stackFrame{}, errors.New("empty stack")
	}
	r := g.stack[l-1]
	g.stack = g.stack[:l-1]
	return r, nil
}

// popSilent ejects top of stack's stackFrame ignoring errors.
func (g *globalState) popSilent() {
	_, err := g.pop()
	if err != nil {
		fmt.Println(err)
	}
}

// peek returns top stackFrame without removing it from the stack.
func (g *globalState) peek() *stackFrame {
	l := len(g.stack)
	if l == 0 {
		return nil
	}
	return &g.stack[l-1]
}

// ==== Evaluated Objects ====

// newFolder takes a parsed folder YAML object and stores in for later processing.
//
// newFolder will ignore duplicated object based on the specified id.
func newFolder(f *folderSpecYAML) error {
	if _, ok := gState.evaluated.folders.YAMLs[f.Id]; ok {
		return errors.New("duplicated definition of folder")
	}
	gState.evaluated.folders.YAMLs[f.Id] = f
	return nil
}
