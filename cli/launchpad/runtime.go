// Package launchpad file runtime.go contains runtime related support
// for evaluation hierarchy and eval object tracking.
package launchpad

import (
	"fmt"
	"log"
	"strings"
)

type eval struct {
	rmap referenceMap
	org  orgSpecYAML
}

func evaluate(rs []resourcer) *eval {
	e := eval{}
	e.rmap.ids = make(map[string]resourcespecer)
	e.rmap.mappings = make(map[string][]reference)

	for _, r := range rs {
		if err := r.enroll(&e); err != nil {
			log.Fatalln("error validating YAML", err.Error())
		}
	}
	if err := e.rmap.organize(); err != nil {
		log.Fatalln("unable to organize intake YAMLs", err.Error())
	}
	return &e
}

// register records the source for future usage.
// we basically register all the ids here and initialize the reference map.
// TODO more doc
func (e *eval) register(src resourcespecer, dst *referenceYAML) {
	e.rmap.ids[src.refId()] = src

	if dst.TargetType() == Organization { // Initialize Organization
		if err := e.org.initializeByRef(dst); err != nil { // attempt to initialize org
			log.Fatalln(fmt.Sprintf("org already initialized as %s, cannot initialize to %s",
				e.org.Id, dst.TargetId))
		}
		// Org is special that we need to manually register it's ID for them
		// TODO take care of explicit org initialization case
		e.rmap.ids[e.org.refId()] = &e.org
	}

	if _, found := e.rmap.mappings[dst.refId()]; !found {
		e.rmap.mappings[dst.refId()] = []reference{}
	}
	e.rmap.mappings[dst.refId()] = append(e.rmap.mappings[dst.refId()], reference{srcPtr: src, dstRef: dst})
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
// note that during YAML evaluation it is unclear Folder X will be eval
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
	srcPtr resourcespecer
	dstRef *referenceYAML
}

// referenceMap maps referenced target to a list of source references.
type referenceMap struct {
	mappings map[string][]reference
	ids      map[string]resourcespecer
}

// organize loop over all addressable target and ensure parent child relationships.
func (m *referenceMap) organize() error {
	for k, v := range m.mappings {
		target, found := m.ids[k]
		if !found {
			log.Fatalln(fmt.Sprintf("reference %s does not exist", k))
		}
		if err := target.resolveRefs(v); err != nil {
			return err
		}
	}
	return nil
}

func (m *referenceMap) dump(ind int) []string {
	indent := strings.Repeat(" ", ind)
	buff := []string{fmt.Sprintf("%sReference Maps:", indent)}

	for k, v := range m.mappings {
		var srcs []string
		for _, vv := range v {
			srcs = append(srcs, vv.srcPtr.refId())
		}
		buff = append(buff, fmt.Sprintf("%s  * %s <- [%s]", indent, k, strings.Join(srcs, ", ")))
	}
	buff = append(buff, fmt.Sprintf("%sReferencable Targets:", indent))
	for k := range m.ids {
		buff = append(buff, fmt.Sprintf("%s  - %s", indent, k))
	}
	return buff
}

func (m *referenceMap) String() string {
	return strings.Join(m.dump(0), "\n")
}
