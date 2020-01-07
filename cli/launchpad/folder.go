package launchpad

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

const (
	folderNameMin = 3
	folderNameMax = 30
)

// folderSpecYAML defines GCP Folder's spec.
type folderSpecYAML struct { // Inner mappings
	Id             string                 `yaml:"id"`
	DisplayName    string                 `yaml:"displayName"`
	ParentRef      referenceYAML          `yaml:"parentRef"`
	SubFolderSpecs []*folderSpecYAML      `yaml:"folders"`
	Undefined      map[string]interface{} `yaml:",inline"` // Catch-all for unintended behavior
}

// folders represents a list of folders.
type folders []*folderYAML

// add appends a folder into the folder list if it does not exist already.
func (fs *folders) add(newF *folderYAML) {
	for _, f := range *fs {
		if f.Spec.Id == newF.Spec.Id { // silently ignore already existing folder
			return
		}
	}
	*fs = append(*fs, newF)
}

func (fs folders) sortedCopy() folders {
	buff := make(folders, len(fs)) // sorted subFolders
	copy(buff, fs)
	sort.SliceStable(buff, func(i, j int) bool { return buff[i].Spec.Id < buff[j].Spec.Id })
	return buff
}

// folderYAML is a GCP Folder YAML representation.
type folderYAML struct {
	headerYAML `yaml:",inline"`
	Spec       folderSpecYAML `yaml:"spec"`
	subFolders folders        // subFolders is a validated sub directories.
}

// resId returns an internal referencable id.
func (f *folderYAML) resId() string { return fmt.Sprintf("%s.%s", Folder, f.Spec.Id) }

// String implements Stringer and generates a string representation.
func (f *folderYAML) String() string { return strings.Join(f.dump(0), "\n") }

// validate ensures input YAML fields are correct.
//
// validate also populates subFolders.
func (f *folderYAML) validate() error {
	if f.Spec.Id == "" {
		return errMissingRequiredField
	}
	switch f.Spec.ParentRef.TargetTypeStr { // Validate Supported Parents
	case Organization.String(), Folder.String():
	default:
		log.Printf("fatal: unsupported parent '%s' type for Folder\n", f.Spec.ParentRef.TargetTypeStr)
		return errInvalidParent
	}

	if ind := tfNameRegex.FindStringIndex(f.Spec.Id); ind == nil {
		log.Printf("GCP Folder [%s] ID does not conform to Terraform standard", f.Spec.DisplayName)
		return errValidationFailed
	}

	if len(f.Spec.DisplayName) < folderNameMin || len(f.Spec.DisplayName) > folderNameMax {
		log.Printf("GCP Folder Name [%s] needs to be between %d and %d", f.Spec.DisplayName, folderNameMin, folderNameMax)
		return errValidationFailed
	}

	f.subFolders = newSubFoldersBySpecs(f.Spec.SubFolderSpecs, Folder, f.Spec.Id)
	return nil
}

// addToOrg adds the folder into the assembled organization.
//
// addToOrg also recursively add folder's subFolders into the org.
func (f *folderYAML) addToOrg(ao *assembledOrg) error {
	if err := ao.registerResource(f, &f.Spec.ParentRef); err != nil {
		return err
	}

	for _, sf := range f.subFolders { // Recursively enroll sub-folders
		if err := sf.addToOrg(ao); err != nil {
			return err
		}
	}
	return nil
}

// resolveReferences processes references to folder.
//
// resolveReferences takes reference from folder as a subFolder of this folder.
func (f *folderYAML) resolveReferences(refs []resourceHandler) error {
	for _, ref := range refs {
		switch r := ref.(type) {
		case *folderYAML:
			if f.Spec.Id != r.Spec.ParentRef.TargetId { // caller should already ensure this once
				log.Printf("fatail: mismatch parent id %s %s", f.resId(), r.Spec.ParentRef.refId())
				return errInvalidParent
			}
			f.subFolders.add(r)
		default:
			log.Printf("fatal: invalid %s parent for %s\n", f.resId(), r.resId())
			return errInvalidInput
		}
	}
	return nil
}

// dump generates debug string slices representation.
func (f *folderYAML) dump(ind int) []string {
	indent := strings.Repeat(" ", ind)
	rep := fmt.Sprintf("%s+ %s.%s (\"%s\") < %s.%s", indent, Folder, f.Spec.Id,
		f.Spec.DisplayName, f.Spec.ParentRef.TargetTypeStr, f.Spec.ParentRef.TargetId)
	buff := []string{rep}

	for _, sf := range f.subFolders.sortedCopy() {
		buff = append(buff, sf.dump(ind+defaultIndentSize)...)
	}
	return buff
}

// newSubFoldersBySpecs initializes folderSpecYAMLs and turn it into a folderYAMLs.
//
// newSubFoldersBySpecs overwrites folderSpecYAML's parent field if parentId is provided.
func newSubFoldersBySpecs(sfYAMLs []*folderSpecYAML, parentType crdKind, parentId string) []*folderYAML {
	var sfs []*folderYAML

	for _, sfYAML := range sfYAMLs {
		sf := folderYAML{Spec: *sfYAML}
		if parentId != "" { // overwrite parents setting
			sf.Spec.ParentRef.TargetTypeStr = parentType.String()
			sf.Spec.ParentRef.TargetId = parentId
		}
		sf.APIVersion = apiCFTv1alpha1
		sf.KindStr = Folder.String()

		sfs = append(sfs, &sf)
	}
	return sfs
}
