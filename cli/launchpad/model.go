// Package launchpad file model.go contains all supported CustomResourceDefinition (CRD)
//
// Every CRD should have a `{kind}YAML` to denote the full CRD representation,
// and a `{kind}SpecYAML` to denote the spec map inside the CRD
//
// Each `{kind}SpecYAML` should also implement stackable interface (synonymous to yaml.Unmarshaler)
// to allow a stack like evaluation
package launchpad

import (
	"errors"
	"fmt"
	"log"
)

// configYAML represents fully qualified CRD that can be of any supported Kind.
type configYAML struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       crdKind     `yaml:"kind"`
	Spec       interface{} `yaml:"spec"` // Placeholder, yielding for further evaluation
	//Metadata   Metadata               `yaml:"metadata"`
	Undefined map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior

	// Original YAML for backup purpose
	rawYAML string
}

// Abstraction of common YAML attributes
type commonConfigYAML configYAML

// UnmarshalYAML evaluates the common attributes and dynamically parse the specified Kind.
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
type parentRefYAML struct {
	ParentType crdKind `yaml:"type"`
	ParentId   string  `yaml:"id"`
}

// ==== CloudFoundation ====

// cloudFoundationYAML represents CloudFoundation CRD.
type cloudFoundationYAML struct {
	commonConfigYAML
	Spec cloudFoundationSpecYAML `yaml:"spec"`
}

// cloudFoundationSpecYAML defines CloudFoundation Kind's spec.
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
func (f *folderSpecYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	top := gState.peek()
	gState.push(KindFolder, f)
	defer gState.popSilent()

	type raw folderSpecYAML
	if err := unmarshal((*raw)(f)); err != nil {
		return err
	}

	// ParentRef can be default value if it is nested under others
	//
	//   kind: Folder
	//   spec:
	//     id: X
	//     folders:
	//       - id: Y
	//
	// During evaluation of Y, ParentRef will not be set, however, stack
	// will retain its owner X and hold its reference
	if f.ParentRef.ParentId == "" {
		if top == nil {
			// ParentRef is not provided, and there are no ownership mappings
			log.Printf("warning, cannot infer parents, output will need to be manually filled in")
		} else {
			switch parent := top.stackPtr.(type) {
			case *folderSpecYAML:
				f.ParentRef.ParentType = KindFolder
				f.ParentRef.ParentId = parent.Id
			case *orgSpecYAML:
				f.ParentRef.ParentType = KindOrganization
				f.ParentRef.ParentId = parent.Id
			}
		}
	}
	// TODO (wengm) check undefined

	// TODO validate folder ID format
	// TODO validate folder displayName format
	if err := newFolder(f); err != nil {
		log.Println("warning: ignoring duplicated definition of folder", f.Id)
	}
	return nil
}
