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
	KindFolder crdKind = "Folder"
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
	case KindFolder:
		err = unmarshal(&folderYAML{})
	default:
		return errors.New(fmt.Sprintf("Kind %s not implemented", c.Kind))
	}
	return err
}

type parentRefYAML struct {
	ParentType string `yaml:"type"`
	ParentId   string `yaml:"id"`
}

// ==== Folder ====
type folderYAML struct {
	commonConfigYAML
	Spec struct {
		Id          string                 `yaml:"id"`
		DisplayName string                 `yaml:"displayName"`
		ParentRef   parentRefYAML          `yaml:"parentRef"`
		Folders     []*folderYAML          `yaml:"folders"`
		Undefined   map[string]interface{} `yaml:",inline"` // Catch-all for untended behavior1
	} `yaml:"spec"`
}

func (f *folderYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type raw folderYAML
	if err := unmarshal((*raw)(f)); err != nil {
		return err
	}
	//checkUndefined("folder", f.Undefined)

	// TODO validate folder ID format
	// TODO validate folder displayName format

	if err := newFolder(f); err != nil {
		log.Println("warning: ignoring duplicated definition of folder", f.Spec.Id)
	}
	return nil
}