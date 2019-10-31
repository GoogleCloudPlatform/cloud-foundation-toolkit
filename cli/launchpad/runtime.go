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
	"log"
)

var (
	errDuplicatedDefinition = errors.New("duplicated definition")
	errOrgIdConflict        = errors.New("conflicted organization definition")
	errUndefinedReference   = errors.New("undefined references")
	errUnsupportedReference = errors.New("unsupported reference type")
)

// globalState stores metadata such as stack and evaluated objects to facilitate output generation.
type globalState struct {
	stack           []stackFrame
	outputDirectory string
	outputFlavor    outputFlavor
	references      map[string][]reference
	evaluated       struct {
		orgId   string
		folders folders
	}
}

// init resets globalState to initialized state.
func (g *globalState) init() {
	g.stack = []stackFrame{}
	g.outputDirectory = ""
	g.outputFlavor = ""
	g.references = make(map[string][]reference)
	g.evaluated.orgId = ""
	g.evaluated.folders.YAMLs = make(map[string]*folderSpecYAML)
}

// reference is an explicit reference specified by the CRD.
//
// An explicit reference can be of any type from the CRD to reference
// to another resource. Source is the origin of the reference, and destination
// is the targeted resource, which may or may not exist.
type reference struct {
	srcType crdKind
	srcPtr  stackable
	dstType crdKind
	dstId   string
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

// ==== CRD References ====

// storeReference takes a given reference and preserve it for validation later.
//
// Some CRD will specify explicit references, such as which organization/folder current project belongs.
// In the case where user specifies a resource ID that does not exist, storeReference and checkReferences
// works in conjunction to help prevent unknown reference error.
func (g *globalState) storeReference(origin stackFrame, r *parentRefYAML) {
	refId := string(r.ParentType) + "." + r.ParentId
	ref := reference{origin.stackType, origin.stackPtr, r.ParentType, r.ParentId}
	g.references[refId] = append(g.references[refId], ref)
}

// checkReferences validates all the explicit reference processed from CRDs and report error on undefined references.
func (g *globalState) checkReferences() error {
	for _, refs := range g.references {
		for _, ref := range refs {
			switch ref.dstType {
			case KindFolder:
				if refFolder, ok := gState.evaluated.folders.YAMLs[ref.dstId]; !ok {
					log.Printf("Folder reference [%s] undefined\n", ref.dstId)
					return errUndefinedReference
				} else {
					if err := refFolder.resolveReference(ref); err != nil {
						return err
					}
				}
			case KindOrganization:
				if gState.evaluated.orgId != ref.dstId {
					log.Printf("Organization reference [%s] conflicted\n", ref.dstId)
					return errOrgIdConflict
				}
			default:
				log.Println("Unknown Reference Type", ref.dstType)
				return errUnsupportedReference
			}
		}
	}
	return nil
}

// ==== Evaluated Resources ====

// setOrg takes a organization id and verify it is not conflict with previous definitions.
//
// Multiple CRDs can be processed at a given time, and Launchpad only allows generation
// of one organization. Therefore, all organization IDs in CRDs must be the same.
func (g *globalState) setOrg(orgId string) error {
	if g.evaluated.orgId == "" {
		g.evaluated.orgId = orgId
		return nil
	}
	if g.evaluated.orgId == orgId {
		return nil
	}
	return errOrgIdConflict
}

// storeFolder takes a parsed folder YAML object and stores in for later processing.
//
// storeFolder will ignore duplicated object based on the specified id.
func (g *globalState) storeFolder(f *folderSpecYAML) error {
	if _, ok := gState.evaluated.folders.YAMLs[f.Id]; ok {
		return errDuplicatedDefinition
	}
	gState.evaluated.folders.YAMLs[f.Id] = f
	return nil
}
