package launchpad

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	folderNameMin = 3
	folderNameMax = 30
)

// folderYAML is a GCP Folder YAML representation.
type folderYAML struct {
	headerYAML `yaml:",inline"`
	Spec       folderSpecYAML `yaml:"spec"`
}

func (f *folderYAML) validate() error      { return f.Spec.validate() }
func (f *folderYAML) enroll(e *eval) error { return f.Spec.enroll(e) }
func (f *folderYAML) String() string       { return strings.Join(f.Spec.dump(0), "\n") }

// folderSpecYAML defines GCP Folder's spec.
type folderSpecYAML struct { // Inner mappings
	Id          string                 `yaml:"id"`
	DisplayName string                 `yaml:"displayName"`
	ParentRef   referenceYAML          `yaml:"parentRef"`
	Folders     folderSpecYAMLs        `yaml:"folders"`
	Undefined   map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior
}

func (f *folderSpecYAML) refId() string { return fmt.Sprintf("%s.%s", Folder, f.Id) }
func (f *folderSpecYAML) dump(ind int) []string {
	indent := strings.Repeat(" ", ind)
	rep := fmt.Sprintf("%s+ %s.%s (\"%s\") <- (%s.%s)", indent, Folder, f.Id,
		f.DisplayName, f.ParentRef.TargetTypeStr, f.ParentRef.TargetId)
	buff := []string{rep}
	for _, sf := range f.Folders {
		buff = append(buff, sf.dump(ind+indentSize)...)
	}
	return buff
}

func (f *folderSpecYAML) validate() error {
	switch f.ParentRef.TargetType() { // Validate Supported Parents
	case Organization, Folder:
	default:
		return errors.New(fmt.Sprintf("unsupported parent '%s' type for Folder", f.ParentRef.TargetTypeStr))
	}

	// TODO validate misc
	if ind := tfNameRegex.FindStringIndex(f.Id); ind == nil {
		log.Printf("GCP Folder [%s] ID does not conform to Terraform standard", f.DisplayName)
		return errValidationFailed
	}

	if len(f.DisplayName) < folderNameMin || len(f.DisplayName) > folderNameMax {
		log.Printf("GCP Folder Name [%s] needs to be between %d and %d", f.DisplayName, folderNameMin, folderNameMax)
		return errValidationFailed
	}

	for _, sf := range f.Folders {
		// Setting & Over writing parent reference for sub-folders
		sf.ParentRef.TargetTypeStr = Folder.String()
		sf.ParentRef.TargetId = f.Id

		if err := sf.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (f *folderSpecYAML) enroll(e *eval) error {
	e.register(f, &f.ParentRef)
	for _, sf := range f.Folders { // Recursively enroll sub-folders
		if err := sf.enroll(e); err != nil {
			return err
		}
	}
	return nil
}

func (f *folderSpecYAML) resolveRefs(refs []reference) error {
	for _, r := range refs {
		sf, ok := r.srcPtr.(*folderSpecYAML)
		if !ok {
			return errors.New("unable to add non folder as a sub-folder")
		}
		f.Folders.add(sf)
	}
	return nil
}

type folderSpecYAMLs []*folderSpecYAML

func (fs *folderSpecYAMLs) add(newF *folderSpecYAML) {
	for _, f := range *fs {
		if f.Id == newF.Id { // silently ignore already existing folder
			return
		}
	}
	*fs = append(*fs, newF)
}
