package launchpad

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const indentSize = 2

var (
	errValidationFailed = errors.New("validation failed")
	tfNameRegex         = regexp.MustCompile("^[a-zA-Z][a-zA-Z\\d\\-\\_]*$")
)

//
type resourcer interface {
	// validate ensures the parsed YAML has validate parameters.
	validate() error
	// kind returns the validated crdKind of the resource.
	kind() crdKind
	// enroll adds the resource into final eval mapping.
	enroll(*eval) error
}

type resourcespecer interface {
	// refId defines the internal referencable id.
	refId() string
	// validate ensure the parsed YAML has validate parameters.
	validate() error
	// dump generates printable information on the resource.
	dump(int) []string
	// resolveRefs adds references to this resource as referenced target.
	resolveRefs([]reference) error
}

// crdKind is the CustomResourceDefinition (CRD) which is indicated YAML's Kind value.
type crdKind int

const (
	CloudFoundation crdKind = iota
	Folder
	Organization
)

func (k crdKind) String() string {
	return []string{"CloudFoundation", "Folder", "Organization"}[k]
}

// newCRDKind parses string formatted crdKind and convert to internal format.
//
// Unsupported format given will terminate the application.
func newCRDKind(crdKindStr string) crdKind {
	switch strings.ToLower(crdKindStr) {
	case "cloudfoundation":
		return CloudFoundation
	case "folder":
		return Folder
	case "organization":
		return Organization
	default:
		log.Fatalln("Unsupported CustomResourceDefinition", crdKindStr)
	}
	return -1
}

var supportedVersion = map[string]supportedKind{
	"cft.dev/v1alpha1": v1alpha1SupportedKind,
}

type supportedKind map[crdKind]func() resourcer

var v1alpha1SupportedKind = supportedKind{
	Folder:       func() resourcer { return &folderYAML{} },
	Organization: func() resourcer { return &orgYAML{} },
}

type headerYAML struct {
	APIVersion string `yaml:"apiVersion"`
	KindStr    string `yaml:"kind"`
}

func (h *headerYAML) kind() crdKind {
	return newCRDKind(h.KindStr)
}

// referenceYAML represents a reference to another object within a CRD.
//
// Among different types of CRDs, it is common to have parent-children relationship, ex: Projects belong to
// a folder, Network belong to a project. referenceYAML is a relationship identifier between these objects.
//
// With explicit definition of ParentRef is possible
type referenceYAML struct {
	TargetTypeStr string `yaml:"type"`
	TargetId      string `yaml:"id"`
}

func (r *referenceYAML) TargetType() crdKind {
	return newCRDKind(r.TargetTypeStr)
}

func (r *referenceYAML) refId() string {
	return fmt.Sprintf("%s.%s", r.TargetType(), r.TargetId)
}
