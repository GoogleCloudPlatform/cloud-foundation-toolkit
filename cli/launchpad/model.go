// Package launchpad file model.go contains all supported CustomResourceDefinition (CRD).
//
// Every CRD should have a `{kind}YAML` to denote the full CRD representation,
// and a `{kind}SpecYAML` to denote the spec map inside the CRD
//
// Each `{kind}SpecYAML` should also implement stackable interface (synonymous to yaml.Unmarshaler)
// to allow a stack like evaluation.
package launchpad

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
	"strings"
)

var (
	errValidationFailed = errors.New("validation failed")
	tfNameRegex         = regexp.MustCompile("^[a-zA-Z][a-zA-Z\\d\\-\\_]*$")
)

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

// genericYAML represents a fully qualified CRD that can be of any Kind.
//
// All CRDs are expected to have the following properties.
type genericYAML struct {
	APIVersion string                 `yaml:"apiVersion"`
	KindStr    string                 `yaml:"kind"`
	Spec       interface{}            `yaml:"spec"`    // Placeholder, yielding for further evaluation
	Undefined  map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior
}

// Kind returns a validated crdKind type based on YAML Kind field.
func (g *genericYAML) Kind() crdKind {
	return newCRDKind(g.KindStr)
}

// UnmarshalYAML evaluates common CRD attributes and dynamically parse the CRDs based on CRD Kind.
func (g *genericYAML) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type raw genericYAML
	if err := unmarshal((*raw)(g)); err != nil {
		return err
	}
	switch g.Kind() {
	case CloudFoundation:
		return unmarshal(&struct {
			Spec cloudFoundationSpecYAML `yaml:"spec"`
		}{})
	case Organization:
		return unmarshal(&struct {
			Spec orgSpecYAML `yaml:"spec"`
		}{})
	case Folder:
		return unmarshal(&struct {
			Spec folderSpecYAML `yaml:"spec"`
		}{})
	default:
		return fmt.Errorf("crd %s not supported", g.KindStr)
	}
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

// UnmarshalYAML evaluates referenceYAML.
//
// UnmarshalYAML will have side effect to set organization ID if reference type is Organization. And will store the
// referenceMap as this represents an explicit reference.
func (r *referenceYAML) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	src := gState.evalStack.peek() // retrieve reference source
	type raw referenceYAML
	if err := unmarshal((*raw)(r)); err != nil {
		return err
	}

	gState.referenceMap.putExplicit(src, r) // retain explicit reference to check if reference exist
	return nil
}

// ==== CloudFoundation ====
// cloudFoundationSpecYAML defines CloudFoundation Kind's spec.
//
// A CloudFoundation CRD can represent a birds-eye view of the entire organization's GCP environment.
// The CRD can be further broken down into small CRDs (ex: Project, Folder, Org) to represent the same.
//
// It is assumed that one CloudFoundation will host at most one Organization.
type cloudFoundationSpecYAML struct {
	Org orgSpecYAML `yaml:"organization"`
}

// UnmarshalYAML evaluates cloudFoundationSpecYAML.
//
// UnmarshalYAML will have side effect to push CloudFoundation onto its stack while nested
// evaluation is ongoing.
func (c *cloudFoundationSpecYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	gState.evalStack.push(c)
	defer gState.evalStack.popSilent()

	type raw cloudFoundationSpecYAML
	if err := unmarshal((*raw)(c)); err != nil {
		return err
	}
	return nil
}

// ==== Organization ====
// orgSpecYAML defines Organization Kind's spec.
type orgSpecYAML struct {
	Id          string          `yaml:"id"`
	DisplayName string          `yaml:"displayName"`
	Folders     folderSpecYAMLs `yaml:"folders"`
}

func (o *orgSpecYAML) String() string {
	return Organization.String() + "." + o.Id
}

// UnmarshalYAML evaluates orgSpecYAML.
//
// UnmarshalYAML will have side effect to push Organization onto its stack while nested
// evaluation is ongoing.
func (o *orgSpecYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	gState.evalStack.push(o)
	defer gState.evalStack.popSilent()

	type raw orgSpecYAML
	if err := unmarshal((*raw)(o)); err != nil {
		return err
	}

	if err := gState.validatedOrg.merge(o); err != nil {
		return err
	}
	return nil
}

//
func (o *orgSpecYAML) merge(newO *orgSpecYAML) error {
	if o.Id != "" && o.Id != newO.Id {
		return errConflictDefinition
	}
	o.Id = newO.Id
	o.DisplayName = newO.DisplayName
	o.Folders.merge(newO.Folders)
	return nil
}

func (o *orgSpecYAML) StringIndent(indentCount int) string {
	var buff strings.Builder
	indent := strings.Repeat(" ", indentCount)
	buff.WriteString(fmt.Sprintf("%sOrganization (%s:\"%s\")\n", indent, o.Id, o.DisplayName))
	for _, f := range o.Folders {
		buff.WriteString(f.StringIndent(indentCount + 2))
	}
	return buff.String()
}

// ==== Folder ====
const (
	folderNameMin = 3
	folderNameMax = 30
)

type folderSpecYAMLs []*folderSpecYAML

func (fs *folderSpecYAMLs) contains(id string) bool {
	for _, f := range *fs {
		if f.Id == id {
			return true
		}
	}
	return false
}

func (fs *folderSpecYAMLs) add(newF *folderSpecYAML) {
	if fs.contains(newF.Id) {
		log.Println("Warning: ignoring duplicated folder definition", newF.Id)
		return
	}
	*fs = append(*fs, newF)
}

func (fs *folderSpecYAMLs) merge(oFs folderSpecYAMLs) {
	for _, of := range oFs {
		fs.add(of)
	}
}

// folderSpecYAML defines Folder Kind's spec.
type folderSpecYAML struct { // Inner mappings
	Id          string                 `yaml:"id"`
	DisplayName string                 `yaml:"displayName"`
	ParentRef   referenceYAML          `yaml:"parentRef"`
	Folders     folderSpecYAMLs        `yaml:"folders"`
	Undefined   map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior
}

func (f *folderSpecYAML) String() string {
	return Folder.String() + "." + f.Id
}

func (f *folderSpecYAML) StringIndent(indentCount int) string {
	var buff strings.Builder
	indent := strings.Repeat(" ", indentCount)
	buff.WriteString(fmt.Sprintf("%sfolder (%s:\"%s\")", indent, f.Id, f.DisplayName))
	for _, subFolder := range f.Folders {
		buff.WriteString(fmt.Sprintf("\n%s", subFolder.StringIndent(indentCount + 2)))
	}
	return buff.String()
}

// UnmarshalYAML evaluates folderSpecYAML.
//
// UnmarshalYAML will have side effect to push Folder onto its stack while nested
// evaluation is ongoing. In addition, storing validated folder onto gState for tracking.
//
// GCP Folder names required to be in between 3 to 30 characters
func (f *folderSpecYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	top := gState.evalStack.peek()
	gState.evalStack.push(f)
	defer gState.evalStack.popSilent()

	type raw folderSpecYAML
	if err := unmarshal((*raw)(f)); err != nil {
		return err
	}

	// TODO (wengm) warn undefined
	if ind := tfNameRegex.FindStringIndex(f.Id); ind == nil {
		log.Printf("GCP Folder [%s] ID does not conform to Terraform standard", f.DisplayName)
		return errValidationFailed
	}

	if len(f.DisplayName) < folderNameMin || len(f.DisplayName) > folderNameMax {
		log.Printf("GCP Folder Name [%s] needs to be between %d and %d", f.DisplayName, folderNameMin, folderNameMax)
		return errValidationFailed
	}

	if top != nil { // Top of stack exists, parent is implicitly referenced
		if f.ParentRef.TargetId != "" {
			log.Println("Warning: folder", f.Id, "contains both explicit and implicit parents, using implicit")
		}
		f.ParentRef = gState.referenceMap.putImplicit(f, top)
	} else { // Implies this YAML doc is a Folder CRD
		if f.ParentRef.TargetId == "" {
			log.Fatalln("folder", f.Id, "does not have a parent defined")
		}
	}
	return nil
}

// resolveReference processes reference to this folder and merge source under it's
// Folders attribute if valid.
//
// An explicit reference specified current folder as parent. Caller can be of any type,
// however, only supported resources can be a children of this folder.
func (f *folderSpecYAML) resolveReference(ref reference, targetPtr yaml.Unmarshaler) error {
	switch parent := targetPtr.(type) {
	case *folderSpecYAML: // sub folder relationship
		parent.Folders.add(f)
	case *orgSpecYAML: // root folder
		parent.Folders.add(f)
	default:
		return errUnsupportedReference
	}
	return nil
}
