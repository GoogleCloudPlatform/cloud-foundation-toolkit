// Package launchpad file runtime.go contains runtime related support for
// evaluation hierarchy and evaluated object tracking.
//
// In the global state, a evaluation stack is needed to support YAML evaluation. During evaluation, a
// depth-first search approach is employed. For example,
//
//   kind: Folder
//   spec:
//     id: X
//     folders:
//       - id: Y
//
// folder Y will be evaluated first, then folder X. In this case, a stack can help infer parent's attribute
// during the evaluation of folder Y.
package launchpad

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

// globalState stores metadata such as stack and evaluated objects to facilitate output generation.
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
// stackable is synonymous to yaml.Unmarshaler since every CRD's Spec is the target of a stack reference, and
// all Specs implements Unmarshaler to validate input.
type stackable yaml.Unmarshaler

// stackFrame defines a single evaluation hierarchy.
type stackFrame struct {
	stackType crdKind
	stackPtr  stackable
}

// push pushes a new stackFrame onto the current stack.
func (g *globalState) push(stackType crdKind, stackPtr stackable) {
	g.stack = append(g.stack, stackFrame{
		stackType: stackType,
		stackPtr:  stackPtr,
	})
}

// pop ejects top of stack's stackFrame.
func (g *globalState) pop() (*stackFrame, error) {
	l := len(g.stack)
	if l == 0 {
		return nil, errors.New("empty stack")
	}
	r := g.stack[l-1]
	g.stack = g.stack[:l-1]
	return &r, nil
}

// popSilent ejects top of stack's stackFrame ignoring errors.
func (g *globalState) popSilent() {
	if _, err := g.pop(); err != nil {
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
