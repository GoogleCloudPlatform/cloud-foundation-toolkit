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

var errUndefinedReference = errors.New("undefined reference")

type resource struct {
	yaml   resourceHandler
	inRefs []resourceHandler
}

type resourceMap map[string]*resource

func (rm resourceMap) getInit(rId string, yaml resourceHandler) *resource {
	res, found := rm[rId]
	if !found { // initialize resource was not encountered before
		res = &resource{yaml: yaml} // yaml is possible be nil
		rm[rId] = res
	}
	if yaml != nil {
		res.yaml = yaml
	}
	return res
}

func (rm resourceMap) addRef(dst string, yaml resourceHandler) {
	dstResource := rm.getInit(dst, nil) // initialize resource before encountering
	dstResource.inRefs = append(dstResource.inRefs, yaml)
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
			log.Fatalln("error validating YAML", err.Error())
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
	_ = ao.resourceMap.getInit(src.resId(), src)

	if dst == nil { // no outgoing reference from src
		return nil
	}

	if dst.TargetType() == Organization { // Initialize Organization
		if err := ao.org.initializeByRef(dst); err != nil { // attempt to initialize org
			return err
		}
		// org is special that we need to manually register it's ID for them
		ao.resourceMap.getInit(ao.org.resId(), &ao.org)
	}

	// update referenceTracker for references from src
	ao.resourceMap.addRef(dst.resId(), src)
	return nil
}

// resolveReferences loops through resourceMap to link up resource to sub resources.
func (ao *assembledOrg) resolveReferences() error {
	// TODO dont do resId do resId
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
