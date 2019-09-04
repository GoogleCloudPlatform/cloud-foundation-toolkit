package launchpad

import (
	"errors"
	"fmt"
)

type outputFlavor string

func conv(f string) outputFlavor {
	return outputFlavor(f)
}

const (
	outDm outputFlavor = "dm"
	outTf outputFlavor = "tf"
)

var gState globalState

func init() {
	gState.evaluated.folders.YAMLs = make(map[string]*folderSpecYAML)
}

type globalState struct {
	// A simple stack implementation
	stack []stackFrame

	outputDirectory string
	outputFlavor    outputFlavor
	evaluated       struct {
		folders folders
	}
}

// ==== Stack Implementation ====
type stackable interface {
	UnmarshalYAML(unmarshal func(interface{}) error) error
}
type stackFrame struct {
	stackType crdKind
	stackPtr  stackable
}

func (g *globalState) push(stackType crdKind, stackPtr stackable) {
	g.stack = append(g.stack, stackFrame{
		stackType: stackType,
		stackPtr:  stackPtr,
	})
}

func (g *globalState) pop() (stackFrame, error) {
	l := len(g.stack)
	if l == 0 {
		return stackFrame{}, errors.New("empty stack")
	}
	r := g.stack[l-1]
	g.stack = g.stack[:l-1]
	return r, nil
}

func (g *globalState) popSilent() {
	_, err := g.pop()
	if err != nil {
		fmt.Println(err)
	}
}

// Lookup the top most recent stack frame
func (g *globalState) peek() *stackFrame {
	l := len(g.stack)
	if l == 0 {
		return nil
	}
	return &g.stack[l-1]
}

// ==== Constructors for runtime state ====
func newFolder(f *folderSpecYAML) error {
	if _, ok := gState.evaluated.folders.YAMLs[f.Id]; ok {
		return errors.New("duplicated definition of folder")
	}
	gState.evaluated.folders.YAMLs[f.Id] = f
	return nil
}
