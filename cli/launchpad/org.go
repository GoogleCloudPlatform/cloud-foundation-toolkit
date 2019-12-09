package launchpad

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

var errConflictId = errors.New("unable to initialize organization to a different id")

// orgSpecYAML defines an Organization's Spec.
type orgSpecYAML struct {
	Id             string            `yaml:"id"`          // GCP organization id.
	DisplayName    string            `yaml:"displayName"` // Optional field to denote GCP organization name.
	SubFolderSpecs []*folderSpecYAML `yaml:"folders"`
}

// orgYAML represents a GCP organization.
type orgYAML struct {
	headerYAML `yaml:",inline"`
	Spec       orgSpecYAML `yaml:"spec"`
	subFolders folders     // subFolder represents validated sub directories.
}

// refId returns an internal referencable id.
func (o *orgYAML) refId() string { return fmt.Sprintf("%s.%s", Organization, o.Spec.Id) }

// String implements Stringer and generates a string representation.
func (o *orgYAML) String() string { return strings.Join(o.dump(0), "\n") }

// validate ensures input YAML fields are correct.
//
// validate also populates subFolders.
func (o *orgYAML) validate() error {
	// TODO validate ORG spec

	o.subFolders = newSubFoldersBySpecs(o.Spec.SubFolderSpecs, Organization, o.Spec.Id)
	return nil
}

// addToOrg adds the organization into the assembled organization.
//
// addToOrg also recursively add organization's subFolders into the org.
func (o *orgYAML) addToOrg(ao *assembledOrg) error {
	if err := ao.registerResource(o, nil); err != nil {
		return err
	}

	for _, sf := range o.subFolders { // Recursively enroll sub-folders
		if err := sf.addToOrg(ao); err != nil {
			return err
		}
	}
	return nil
}

// resolveReferences processes references to organization.
//
// resolveReferences takes reference from folder as a subFolder of this organization.
func (o *orgYAML) resolveReferences(refs []resourceHandler) error {
	for _, ref := range refs {
		switch r := ref.(type) {
		case *folderYAML:
			o.subFolders.add(r)
		default:
			return errors.New("unable to process reference from resource")
		}
	}
	return nil
}

// initializeByRef initializes an organization through another resource's reference.
func (o *orgYAML) initializeByRef(ref *referenceYAML) error {
	if o.Spec.Id != "" && o.Spec.Id != ref.TargetId {
		log.Printf("fatal: org already initialized to %s, cannot reinitialize to %s\n", o.Spec.Id, ref.TargetId)
		return errConflictId
	} else if o.Spec.Id == "" && ref.TargetId == "" {
		log.Printf("fatal: trying to initialize org with empty Id\n")
		return errors.New("unset org id")
	}
	o.Spec.Id = ref.TargetId
	return nil
}

// dump generates debug string slices representation.
func (o *orgYAML) dump(ind int) []string {
	indent := strings.Repeat(" ", ind)
	rep := fmt.Sprintf("%s%s.%s (\"%s\")", indent, Organization, o.Spec.Id, o.Spec.DisplayName)
	buff := []string{rep}

	for _, sf := range o.subFolders {
		buff = append(buff, sf.dump(ind+defaultIndentSize)...)
	}
	return buff
}
