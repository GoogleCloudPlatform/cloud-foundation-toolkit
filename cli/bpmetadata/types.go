package bpmetadata

import "sigs.k8s.io/kustomize/kyaml/yaml"

// BlueprintMetadata defines the overall structure for blueprint metadata details
type BlueprintMetadata struct {
	Meta yaml.ResourceMeta `json:",inline" yaml:",inline"`
	Spec *BlueprintMetadataSpec
}

// BlueprintMetadataSpec defines the spec portion of the blueprint metadata
type BlueprintMetadataSpec struct {
	Info         BlueprintInfo
	Content      BlueprintContent
	Interfaces   BlueprintInterface
	Requirements BlueprintRequirements
}

// BlueprintInfo defines informational detail for the blueprint
type BlueprintInfo struct {
	Title         string
	Source        *BlueprintRepoDetail
	Version       string
	ActuationTool BlueprintActuationTool `json:"actuationTool" yaml:"actuationTool"`
	Description   *BlueprintDescription  `json:",omitempty" yaml:",omitempty"`
	Icon          string                 `json:",omitempty" yaml:",omitempty"`
}

// BlueprintContent defines the detail for blueprint related content such as
// related documentation, diagrams, examples etc.
type BlueprintContent struct {
	Diagrams      []BlueprintDiagram
	Documentation []BlueprintListContent `json:",omitempty" yaml:",omitempty"`
	SubBlueprints []BlueprintMiscContent `json:"subBlueprints,omitempty" yaml:"subBlueprints,omitempty"`
	Examples      []BlueprintMiscContent
}

// BlueprintInterface the input and output variables for the blueprint
type BlueprintInterface struct {
	Variables      []BlueprintVariable
	VariableGroups []BlueprintVariableGroup `json:"variableGroups" yaml:"variableGroups"`
	Outputs        []BlueprintOutput
}

// BlueprintRequirements defines the roles required and the assocaited services
// that need to be enabled to provision blueprint resources
type BlueprintRequirements struct {
	Roles    []BlueprintRoles
	Services []string
}

type BlueprintRepoDetail struct {
	Repo       string
	SourceType string `json:"sourceType" yaml:"sourceType"`
}

type BlueprintActuationTool struct {
	Flavor  string `json:"type" yaml:"type"`
	Version string
}

type BlueprintDescription struct {
	Tagline   string `json:",omitempty" yaml:",omitempty"`
	Detailed  string `json:",omitempty" yaml:",omitempty"`
	PreDeploy string `json:"preDeploy,omitempty" yaml:"preDeploy,omitempty"`
}

type BlueprintMiscContent struct {
	Name     string
	Location string
}

type BlueprintDiagram struct {
	Name        string
	AltText     string `json:"altText,omitempty" yaml:"altText,omitempty"`
	Description string `json:",omitempty" yaml:",omitempty"`
}

type BlueprintListContent struct {
	Title string
	Url   string `json:",omitempty" yaml:",omitempty"`
}

type BlueprintVariable struct {
	Name        string
	Description string
	VarType     string `yaml:"type"`
	Default     interface{}
	Required    bool
}

type BlueprintVariableGroup struct {
	Name        string
	Description string
	Variables   []string
}

type BlueprintOutput struct {
	Name        string
	Description string
}

type BlueprintRoles struct {
	Level string
	Roles []string
}
