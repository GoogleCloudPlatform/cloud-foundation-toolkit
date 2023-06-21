package launchpad

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	// defaultIndentSize defines default indent size when calling String and dump.
	defaultIndentSize = 2
	apiCFTv1alpha1    = "cft.dev/v1alpha1"
)

// supportedVersion defines API version and support resource for each version.
var supportedVersion = map[string]map[crdKind]func() resourceHandler{
	apiCFTv1alpha1: v1alpha1SupportedKind,
}

// v1alpha1SupportedKind defines v1alpha1's supported resources.
var v1alpha1SupportedKind = map[crdKind]func() resourceHandler{
	Folder:       func() resourceHandler { return &folderYAML{} },
	Organization: func() resourceHandler { return &orgYAML{} },
}

var (
	errValidationFailed     = errors.New("validation failed")
	errMissingRequiredField = errors.New("missing required field")
	errInvalidParent        = errors.New("invalid parent reference")
	errInvalidInput         = errors.New("invalid input")
	tfNameRegex             = regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d\-\_]*$`)
)

// resourceHandler represents a resource that can be processed by launchpad.
type resourceHandler interface {
	// resId defines the internal referencable id.
	resId() string
	// validate ensures the parsed YAML has validate fields.
	validate() error
	// kind returns the validated crdKind of the resource.
	kind() crdKind
	// addToOrg adds the resource into the assembled organization.
	addToOrg(ao *assembledOrg) error
	// resolveReferences takes action on resources referencing the current resource.
	resolveReferences(refs []resourceHandler) error
}

// crdKind is the CustomResourceDefinition (CRD) which is indicated by YAML Kind value.
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
// Unsupported format given will return -1, caller is expected to handle unknown type.
func newCRDKind(crdKindStr string) crdKind {
	switch strings.ToLower(crdKindStr) {
	case "cloudfoundation":
		return CloudFoundation
	case "folder":
		return Folder
	case "organization":
		return Organization
	default:
		log.Printf("warning: unsupported CustomResourceDefinition %q", crdKindStr)
		return -1
	}
}

// headerYAML defines the common fields all CRD is required to have.
type headerYAML struct {
	APIVersion string `yaml:"apiVersion"`
	KindStr    string `yaml:"kind"`
}

func (h *headerYAML) kind() crdKind { return newCRDKind(h.KindStr) }

// referenceYAML represents an explicit reference to another resource.
//
// It is common to have reference relationship among difference resources. For example,
// parent-children relationship such as project belong to a folder, network belong to a
// project. referenceYAML is the relationship identifier between these resources.
type referenceYAML struct {
	TargetTypeStr string `yaml:"type"`
	TargetId      string `yaml:"id"`
}

func (r *referenceYAML) TargetType() crdKind { return newCRDKind(r.TargetTypeStr) }
func (r *referenceYAML) resId() string       { return fmt.Sprintf("%s.%s", r.TargetType(), r.TargetId) }
