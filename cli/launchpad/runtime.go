// Package launchpad file runtime.go contains runtime related support
// for evaluation hierarchy and evaluated object tracking.
package launchpad

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

var (
	errConflictDefinition   = errors.New("conflicted definition")
	errOrgIdConflict        = errors.New("conflicted organization definition")
	errUndefinedReference   = errors.New("undefined referenceMap")
	errUnsupportedReference = errors.New("unsupported reference type")
)

// globalState stores metadata such as stack and evaluated objects to facilitate output generation.
type globalState struct {
	outputDirectory string
	outputFlavor    outputFlavor
	referenceMap    referenceMap
	evalStack       evalStack
	validatedOrg    orgSpecYAML
}

// init resets globalState to initialized state.
func (g *globalState) init() {
	g.evalStack = []yaml.Unmarshaler{}
	g.outputDirectory = ""
	g.outputFlavor = -1
	g.referenceMap.mappings = make(map[string][]reference)
	g.referenceMap.ids = make(map[string]yaml.Unmarshaler)
}

// dump prints out the current globalState to standard output.
func (g *globalState) dump() {
	var buff strings.Builder
	buff.WriteString("Current State:\n==================================\n")
	buff.WriteString(fmt.Sprintf("Output Directory: %s\n", g.outputDirectory))
	buff.WriteString(fmt.Sprintf("Validated Org: %s\n", gState.validatedOrg.StringIndent(0)))
	buff.WriteString(gState.referenceMap.String())
	buff.WriteString("Evaluated Objects: ")
	for k := range g.referenceMap.ids {
		buff.WriteString(fmt.Sprintf("%s ", k))
	}
	buff.WriteString("\n")

	log.Print(buff.String())
}

// ==== Evaluation Stack ====
// An evaluation stack is needed to support YAML evaluation. During
// unmarshal process, a depth-first search approach is employed by
// the go-yaml library. That is,
//
//   kind: Folder
//   spec:
//     targetId: X
//     folderSpecYAMLs:
//       - targetId: Y
//
// folder Y will be fully evaluated first, then folder X. A stack
// can help to infer folder X's attributes during folder Y's evaluation.
//
// As all nested attributes eventually is a YAML struct and all
// meaningful CRDs should implement yaml.Unmarshaler, stack will be
// operated as such.

// evalStack is an evaluation tracker fo facilitate YAML evaluation order.
type evalStack []yaml.Unmarshaler

// push pushes a new layer onto the current evaluation stack.
func (s *evalStack) push(stackPtr yaml.Unmarshaler) {
	*s = append(*s, stackPtr)
}

// pop ejects top of stack's reference.
func (s *evalStack) pop() (yaml.Unmarshaler, error) {
	l := len(*s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}
	r := (*s)[l-1]
	*s = (*s)[:l-1]
	return r, nil
}

// popSilent ejects top of stack's reference ignoring errors.
func (s *evalStack) popSilent() {
	if _, err := s.pop(); err != nil {
		fmt.Println(err)
	}
}

// peek returns top reference without removing it from the stack.
func (s *evalStack) peek() yaml.Unmarshaler {
	l := len(*s)
	if l == 0 {
		return nil
	}
	return (*s)[l-1]
}

// ==== CRD References ====
// Object can be referenced in two distinct ways, implicitly or explicitly.
//
// An implicit reference occurs if an object is nested under another object.
//
//   kind: Folder
//   spec:
//     targetId: X
//     folderSpecYAMLs:
//       - targetId: Y
//
// During the evaluation of folder Y, evalStack can help infer folder X as
// folder Y's parent. Hence, folder Y holds an implicit reference to
// folder X.
//
// An explicit reference occurs if an object specified its type and targetId.
//
//   kind: Folder
//   spec:
//     targetId: Y
//     parentRef:
//       type: Folder
//       targetId: X
//
// During evaluation of folder Y, a reference of Folder X is specified. However,
// note that during YAML evaluation it is unclear Folder X will be evaluated
// or not. Therefore, all explicit references will be stored in a referenceMap
// until all YAMLs are processed then determine if the reference exists.
//
// As reference can have multiple use cases, all YAML definition should
// use `Ref` as suffix for referenced fields. For example, `parentRef` to specify
// parent-child relationship such as organization/folder, folder/folder.

// reference is an explicit reference specified by the CRD.
//
// A reference can be of any type from the CRD to reference to another resource.
// Source is the origin of the reference, and destination
// is the targeted resource, which may or may not exist.
type reference struct {
	srcPtr yaml.Unmarshaler
	dstRef *referenceYAML
}

// targetId generates an unique internal only targetId based on referenced target information.
func (r *reference) targetId() string {
	return r.dstRef.TargetTypeStr + "." + r.dstRef.TargetId
}

// referenceMap maps referenced target to a list of source references.
type referenceMap struct {
	mappings map[string][]reference
	ids      map[string]yaml.Unmarshaler
}

func (m referenceMap) String() string {
	var buff strings.Builder
	buff.WriteString("Reference Map: \n")
	for k, v := range m.mappings {
		buff.WriteString(fmt.Sprintf("  %s <- [", k))
		for _, ptr := range v {
			buff.WriteString(fmt.Sprintf("%s", ptr.srcPtr))
		}
		buff.WriteString("]\n")
	}
	return buff.String()
}

// putExplicit stores an explicit reference to and preserve it for future de-referencing.
func (m referenceMap) putExplicit(src yaml.Unmarshaler, dstRef *referenceYAML) {
	ref := reference{srcPtr: src, dstRef: dstRef}
	m.mappings[ref.targetId()] = append(m.mappings[ref.targetId()], ref)

	if dstRef.TargetType() == Organization {
		if err := gState.setOrgId(dstRef.TargetId); err != nil {
			log.Fatalln(err.Error())
		}
	}

	srcId := ""
	switch s := src.(type) {
	case *folderSpecYAML:
		srcId = s.String()
	case *orgSpecYAML:
		srcId = s.String()
	default:
		log.Fatalln("reference source not supported")
	}
	m.ids[srcId] = src
}

func (m referenceMap) putImplicit(src yaml.Unmarshaler, dst yaml.Unmarshaler) referenceYAML {
	ret := referenceYAML{}
	switch d := dst.(type) {
	case *folderSpecYAML:
		ret.TargetTypeStr = Folder.String()
		ret.TargetId = d.Id
	case *orgSpecYAML:
		ret.TargetTypeStr = Organization.String()
		ret.TargetId = d.Id
	default:
		log.Fatalln("implicit reference destination not supported")
	}
	m.putExplicit(src, &ret)
	return ret
}

// validate checks and links existing references and report error on undefined references.
func (m referenceMap) validate() error {
	for _, refs := range m.mappings {
		for _, ref := range refs {
			targetPtr, ok := m.ids[ref.targetId()]
			if !ok {
				return errUndefinedReference
			}
			switch src := ref.srcPtr.(type) {
			case *folderSpecYAML:
				if err := src.resolveReference(ref, targetPtr); err != nil {
					return err
				}
			case *orgSpecYAML:
				println("TODO not yet implemented")
			default:
				log.Println("Warning: Unknown Reference Type", ref.dstRef.TargetTypeStr)
				return errUnsupportedReference
			}
		}
	}
	return nil
}

// ==== Evaluated Resources ====

// setOrg takes a organization targetId and verify it is not conflict with previous definitions.
//
// Multiple CRDs can be processed at a given time and Launchpad only allows generation
// of one organization. Therefore, all organization IDs in CRDs must be the same.
func (g *globalState) setOrgId(id string) error {
	if g.validatedOrg.Id != "" && g.validatedOrg.Id != id {
		return errOrgIdConflict
	}
	g.validatedOrg.Id = id
	g.referenceMap.ids[Organization.String()+"."+id] = &gState.validatedOrg
	return nil
}
