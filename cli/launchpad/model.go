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
	"log"
	"regexp"
)

const (
	folderNameMin = 3
	folderNameMax = 30
)

var (
	errValidationFailed = errors.New("validation failed")
	tfNameRegex         = regexp.MustCompile("^[a-zA-Z][a-zA-Z\\d\\-\\_]*$")
)

// configYAML represents a fully qualified CRD that can be of any supported Kind.
//
// All CRDs are expected to have the following properties.
type configYAML struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       crdKind     `yaml:"kind"`
	Spec       interface{} `yaml:"spec"` // Placeholder, yielding for further evaluation
	//Metadata   Metadata               `yaml:"metadata"`
	Undefined map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior

	// Original YAML for backup purpose
	rawYAML string
}

// commonConfigYAML alias configYAML for all CRDs to implement common properties.
type commonConfigYAML configYAML

// UnmarshalYAML evaluates common CRD attributes and dynamically parse the CRDs based on Kind field value.
func (c *configYAML) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type raw configYAML
	if err := unmarshal((*raw)(c)); err != nil {
		return err
	}
	switch c.Kind {
	case KindCloudFoundation:
		cf := cloudFoundationYAML{}
		cf.APIVersion = c.APIVersion
		cf.Kind = c.Kind
		err = unmarshal(&cf)
	case KindFolder:
		f := folderYAML{}
		f.APIVersion = c.APIVersion
		f.Kind = c.Kind
		err = unmarshal(&f)
	case KindOrganization:
		o := orgYAML{}
		o.APIVersion = c.APIVersion
		o.Kind = c.Kind
		err = unmarshal(&o)
	default:
		return errors.New(fmt.Sprintf("Kind %s not implemented", c.Kind))
	}
	return err
}

// parentRefYAML represents ownership reference inside a CRD.
//
// Among different types of CRDs, it is common to have parent-children relationship, ex: Projects belong to
// a folder, Network belong to a project. parentRefYAML is a relationship identifier between these objects.
//
// With explicit definition of ParentRef is possible
type parentRefYAML struct {
	ParentType crdKind `yaml:"type"`
	ParentId   string  `yaml:"id"`
}

// UnmarshalYAML evaluates parentRefYAML.
//
// UnmarshalYAML will have side effect to set organization ID if reference type is Organization. And will store the
// references as this represents an explicit reference.
func (p *parentRefYAML) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	top := gState.peek()
	type raw parentRefYAML
	if err := unmarshal((*raw)(p)); err != nil {
		return err
	}
	switch p.ParentType {
	case KindOrganization:
		if err := gState.setOrg(p.ParentId); err != nil {
			return err
		}
	}
	gState.storeReference(*top, p) // retain explicit reference to check if reference exist
	return nil
}

// ==== CloudFoundation ====

// cloudFoundationYAML represents a CloudFoundation CRD.
//
// A CloudFoundation CRD can represent a birds-eye view of the entire organization's GCP environment.
// The CRD can be further broken down into small CRDs (ex: Project, Folder, Org) to represent the same.
type cloudFoundationYAML struct {
	commonConfigYAML
	Spec cloudFoundationSpecYAML `yaml:"spec"`
}

// cloudFoundationSpecYAML defines CloudFoundation Kind's spec.
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
	gState.push(KindCloudFoundation, c)
	defer gState.popSilent()

	type raw cloudFoundationSpecYAML
	if err := unmarshal((*raw)(c)); err != nil {
		return err
	}
	return nil
}

// ==== Organization ====

// orgYAML represents Organization CRD.
//
// orgYAML defines a single GCP Organization
type orgYAML struct {
	commonConfigYAML
	Spec orgSpecYAML `yaml:"spec"`
}

// orgSpecYAML defines Organization Kind's spec.
type orgSpecYAML struct {
	Id          string            `yaml:"id"`
	DisplayName string            `yaml:"displayName"`
	Folders     []*folderSpecYAML `yaml:"folders"`
}

// UnmarshalYAML evaluates orgSpecYAML.
//
// UnmarshalYAML will have side effect to push Organization onto its stack while nested
// evaluation is ongoing.
func (o *orgSpecYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	gState.push(KindOrganization, o)
	defer gState.popSilent()

	type raw orgSpecYAML
	if err := unmarshal((*raw)(o)); err != nil {
		return err
	}
	if err := gState.setOrg(o.Id); err != nil {
		return err
	}
	return nil
}

// ==== Folder ====

// folderYAML represents Folder CRD.
type folderYAML struct {
	commonConfigYAML
	Spec folderSpecYAML `yaml:"spec"`
}

// folderSpecYAML defines Folder Kind's spec.
type folderSpecYAML struct { // Inner mappings
	Id          string                 `yaml:"id"`
	DisplayName string                 `yaml:"displayName"`
	ParentRef   parentRefYAML          `yaml:"parentRef"`
	Folders     []*folderSpecYAML      `yaml:"folders"`
	Undefined   map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior1
}

// UnmarshalYAML evaluates folderSpecYAML.
//
// UnmarshalYAML will have side effect to push Folder onto its stack while nested
// evaluation is ongoing. In addition, storing validated folder onto gState for tracking.
//
// GCP Folder names required to be in between 3 to 30 characters
func (f *folderSpecYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	top := gState.peek()
	gState.push(KindFolder, f)
	defer gState.popSilent()

	type raw folderSpecYAML
	if err := unmarshal((*raw)(f)); err != nil {
		return err
	}

	// Implicit ParentRef can be default value if it is nested under others
	//
	//   kind: Folder
	//   spec:
	//     id: X
	//     folders:
	//       - id: Y
	//
	// During evaluation of Y, ParentRef could be undefined, however, stack
	// will retain owner X reference and infer it as parent
	if f.ParentRef.ParentId == "" {
		if top == nil {
			// ParentRef is not provided, and there are no ownership mappings
			log.Printf("warning, cannot infer parents, output will need to be manually filled in")
		} else {
			switch parent := top.stackPtr.(type) {
			case *folderSpecYAML:
				// implicit reference does not need to state since it is guaranteed that parent is resolved
				f.ParentRef.ParentType = KindFolder
				f.ParentRef.ParentId = parent.Id
			case *orgSpecYAML:
				f.ParentRef.ParentType = KindOrganization
				f.ParentRef.ParentId = parent.Id
				if err := gState.setOrg(parent.Id); err != nil {
					return err
				}
			}
		}
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

	if err := gState.storeFolder(f); err != nil {
		log.Println("warning: ignoring duplicated definition of folder", f.Id)
	}
	return nil
}

// resolveReference processes reference to this folder and merge source under it's
// Folders attribute if valid.
//
// An explicit reference specified current folder as parent. Caller can be of any type,
// however, only supported resources can be a children of this folder.
func (f *folderSpecYAML) resolveReference(ref reference) error {
	switch src := ref.srcPtr.(type) {
	case *folderSpecYAML: // sub folder relationship
		for _, f := range f.Folders {
			if src.Id != f.Id {
				continue
			}
			log.Println("Ignore existing reference within folders", f.Id)
			return nil
		}
		f.Folders = append(f.Folders, src)
	default:
		return errUnsupportedReference
	}
	return nil
}
