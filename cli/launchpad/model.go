// Model contains YAML parsing & validation.
package launchpad

import (
	"errors"
	"fmt"
	"log"
)

// Launchpad Supported CRD type
type crdKind string

// Supported CRD definition
const (
	KindCloudFoundation crdKind = "CloudFoundation"
	KindFolder          crdKind = "Folder"
	KindOrganization    crdKind = "Organization"
)

// ==== Shared Structure ====

// All Configuration YAML should start with CRD style
// This structure also serve as evaluation point for each kind of CRD
type configYAML struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       crdKind     `yaml:"kind"`
	Spec       interface{} `yaml:"spec"` // Placeholder, yielding for further evaluation
	//Metadata   Metadata               `yaml:"metadata"`
	Undefined map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior

	// Original YAML for backup purpose
	rawYAML string
}

// Abstraction of common YAML for others to use
type commonConfigYAML configYAML

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

type parentRefYAML struct {
	ParentType crdKind `yaml:"type"`
	ParentId   string  `yaml:"id"`
}

// ==== CloudFoundation ====
type cloudFoundationYAML struct {
	commonConfigYAML
	Spec cloudFoundationSpecYAML `yaml:"spec"`
}
type cloudFoundationSpecYAML struct {
	Org orgSpecYAML `yaml:"organization"`
}

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
type orgYAML struct { // Dedicated CRD
	commonConfigYAML
	Spec orgSpecYAML `yaml:"spec"`
}
type orgSpecYAML struct {
	Id          string            `yaml:"id"`
	DisplayName string            `yaml:"displayName"`
	Folders     []*folderSpecYAML `yaml:"folders"`
}

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
type folderYAML struct { // Dedicated CRD
	commonConfigYAML
	Spec folderSpecYAML `yaml:"spec"`
}

type folderSpecYAML struct { // Inner mappings
	Id          string                 `yaml:"id"`
	DisplayName string                 `yaml:"displayName"`
	ParentRef   parentRefYAML          `yaml:"parentRef"`
	Folders     []*folderSpecYAML      `yaml:"folders"`
	Undefined   map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior1
}

func (f *folderSpecYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	top := gState.peek()
	gState.push(KindFolder, f)
	defer gState.popSilent()

	type raw folderSpecYAML
	if err := unmarshal((*raw)(f)); err != nil {
		return err
	}

	if f.ParentRef.ParentId == "" {
		// Infer from stack
		if top == nil {
			log.Printf("warning, cannot infer parents, output will need to be manually filled in")
		} else {
			switch parent := top.stackPtr.(type) {
			case *folderSpecYAML:
				f.ParentRef.ParentId = parent.Id
				f.ParentRef.ParentType = KindFolder
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
