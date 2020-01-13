package launchpad

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

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

// resId returns an internal referencable id.
func (o *orgYAML) resId() string { return fmt.Sprintf("%s.%s", Organization, o.Spec.Id) }

// validate ensures input YAML fields are correct.
//
// validate also populates subFolders.
func (o *orgYAML) validate() error {
	if o.Spec.Id == "" {
		return errValidationFailed
	}

	o.subFolders = newSubFoldersBySpecs(o.Spec.SubFolderSpecs, Organization, o.Spec.Id)
	for _, f := range o.subFolders {
		if err := f.validate(); err != nil {
			return err
		}
	}
	return nil
}

// addToOrg adds the organization into the assembled organization.
//
// addToOrg also recursively add organization's subFolders into the org.
func (o *orgYAML) addToOrg(ao *assembledOrg) error {
	// assembledOrg.org could have already been initialized by others via reference, or explicitly
	// need to copy all fields over
	if err := o.mergeFields(&ao.org); err != nil {
		return err
	}
	ao.org = *o // replace finalized org as the current org.

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
			_ = o.subFolders.add(r) // silently ignore existing resource
		default:
			return errors.New("unable to process reference from resource")
		}
	}
	return nil
}

// mergeFields merges all fields from input to current resource.
//
// mergeFields only operates on named resource, ignoring YAML Specs.
// However, future version can consider recursively merging all sub resources through
// additional of mergeFields requirement in resourceHandler.
func (o *orgYAML) mergeFields(oldO *orgYAML) error {
	if oldO.APIVersion != "" {
		o.APIVersion = oldO.APIVersion
	}
	if oldO.Spec.DisplayName != "" {
		o.Spec.DisplayName = oldO.Spec.DisplayName
	}
	o.subFolders.merge(&oldO.subFolders)
	// TODO (FR) recursively merge folderSpecYAML projectSpecYAML ...etc
	// resolveReferences ensures output linkage is valid, hence not a priority as this is a cleanup.
	// downside is {resource}SpecYAML sub-resources are misaligned.
	return nil
}

// dump writes resource's string representation into provided buffer.
func (o *orgYAML) dump(ind int, buff io.Writer) error {
	indent := strings.Repeat(" ", ind)
	_, err := fmt.Fprintf(buff, "%s%s.%s (\"%s\")\n", indent, Organization, o.Spec.Id, o.Spec.DisplayName)
	if err != nil {
		return err
	}

	for _, sf := range o.subFolders {
		err = sf.dump(ind+defaultIndentSize, buff)
		if err != nil {
			return err
		}
	}
	return nil
}
