// Package launchpad file runtime.go contains runtime support for reference
// tracking and output assembled object to represent parsed view of the org
// for later processing.
package launchpad

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

var (
	errUndefinedReference = errors.New("undefined reference")
	errConflictDefinition = errors.New("definition conflict")
	errUnexpectedType     = errors.New("unexpected type")
)

type resource struct {
	yaml   resourceHandler
	inRefs []resourceHandler
}

type resourceMap map[string]*resource

func (rm resourceMap) getInit(rId string, yaml resourceHandler) (*resource, error) {
	res, found := rm[rId]
	if !found { // initialize resource was not encountered before
		res = &resource{yaml: yaml} // yaml is possible be nil
		rm[rId] = res
	}
	if yaml == nil {
		return res, nil
	}
	if res.yaml.kind() == Organization {
		// newer organization definition, pull sub-resources into current
		o, ok := res.yaml.(*orgYAML)
		if !ok {
			return nil, errUnexpectedType
		}
		oNew, ok := yaml.(*orgYAML)
		if !ok {
			return nil, errUnexpectedType
		}
		return res, o.mergeFields(oNew)
	}
	if yaml != res.yaml {
		log.Println("conflicting YAML definition detected on", yaml.resId())
		return nil, errConflictDefinition
	}
	res.yaml = yaml
	return res, nil
}

func (rm resourceMap) sortedResId() []string {
	keys := make([]string, len(rm))
	i := 0
	for k := range rm {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// assembledOrg supports reference tracking and allows resources to be organized into an organization tree.
type assembledOrg struct {
	resourceMap resourceMap // tracks seen resources and references.
	org         orgYAML     // finalized view of all validated resources.
}

// newAssembledOrg creates and initializes an assembledOrg.
func newAssembledOrg() *assembledOrg {
	ao := assembledOrg{}
	ao.resourceMap = make(resourceMap)
	ao.org.headerYAML = headerYAML{APIVersion: apiCFTv1alpha1, KindStr: Organization.String()}
	return &ao
}

// String implements Stringer and generates a string representation.
func (ao *assembledOrg) String() string {
	sb := strings.Builder{}
	err := ao.dump(0, &sb)
	if err != nil {
		panic(err.Error())
	}
	return sb.String()
}

// assembleResourcesToOrg takes in resources and assembles into an organization.
func assembleResourcesToOrg(rs []resourceHandler) *assembledOrg {
	ao := newAssembledOrg()

	// discover resources in a DFS style
	// initialize resource into resourceMap or update references if already exist.
	for _, r := range rs {
		if err := r.addToOrg(ao); err != nil {
			log.Println("error validating YAML", err.Error())
			panic(err.Error())
		}
	}
	// assemble each discovered resource onto a finalized org
	if err := ao.resolveReferences(); err != nil {
		log.Println("unable to resolve referenceTracker between YAML resources:", err.Error())
		panic(err.Error())
	}
	return ao
}

// registerResource registers a resource into resourceMap for later resolution.
//
// registerResource will init/update entry resourceMap with for resId as key. If
// the resource being registered (src) has a reference (dst) to another resource does
// not yet exist, it will also initialize the resourceMap for dst resource.
//
// If there are conflicting resources ID, silently pick the latest.
func (ao *assembledOrg) registerResource(src resourceHandler, dst *referenceYAML) error {
	if _, err := ao.resourceMap.getInit(src.resId(), src); err != nil {
		return err
	}

	if dst == nil { // no outgoing reference from src
		return nil
	}

	// initialize a resource on the resource map, pending future definition on YAML
	var dstYAML resourceHandler
	if dst.TargetType() == Organization {
		if ao.org.Spec.Id == "" {
			ao.org.Spec.Id = dst.TargetId
		} else if ao.org.Spec.Id != dst.TargetId {
			log.Printf("fatal: org is identified as %s, cannot remap to %s\n", ao.org.Spec.Id, dst.TargetId)
			return errConflictDefinition
		}
		dstYAML = &ao.org // allow future resolution to pick up finalized org directly
	}

	dstRes, err := ao.resourceMap.getInit(dst.resId(), dstYAML)
	if err != nil {
		return err
	}
	// update referenceTracker for references from src
	dstRes.inRefs = append(dstRes.inRefs, src)
	return nil
}

// resolveReferences loops through resourceMap to link up resource to sub resources.
func (ao *assembledOrg) resolveReferences() error {
	for resId, res := range ao.resourceMap {
		if res.yaml == nil {
			// an item is initialized but the resourceHandler never provided
			// only happen when this item is initialized via inbound reference(s)
			log.Printf("fatal: reference to %s was not found\n", resId)
			return errUndefinedReference
		}
		// each resource holds its own resolving logic
		if err := res.yaml.resolveReferences(res.inRefs); err != nil {
			return err
		}
	}
	return nil
}

// dump writes resource's string representation into provided buffer.
func (ao *assembledOrg) dump(ind int, buff io.Writer) error {
	indent := strings.Repeat(" ", ind)
	_, err := fmt.Fprintf(buff, "%sResource Map [%d]:\n", indent, len(ao.resourceMap))
	if err != nil {
		return err
	}

	for _, resId := range ao.resourceMap.sortedResId() {
		res := ao.resourceMap[resId]
		var refs []string
		for _, refRes := range res.inRefs {
			refs = append(refs, refRes.resId())
		}
		sort.Strings(refs)

		_, err = fmt.Fprintf(buff, "%s  * %s <- [%s]\n", indent, resId, strings.Join(refs, ", "))
		if err != nil {
			return err
		}
	}
	return ao.org.dump(ind, buff)
}
