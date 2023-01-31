package bpmetadata

import "sigs.k8s.io/kustomize/kyaml/yaml"

// BlueprintMetadata defines the overall structure for blueprint metadata details
type BlueprintMetadata struct {
	Meta yaml.ResourceMeta `json:",inline" yaml:",inline"`
	Spec BlueprintMetadataSpec
}

// BlueprintMetadataSpec defines the spec portion of the blueprint metadata
type BlueprintMetadataSpec struct {
	Info         BlueprintInfo         `json:",inline" yaml:",inline"`
	Content      BlueprintContent      `json:",inline" yaml:",inline"`
	Interfaces   BlueprintInterface    `json:",inline" yaml:",inline"`
	Requirements BlueprintRequirements `json:",inline" yaml:",inline"`
}

// BlueprintInfo defines informational detail for the blueprint
type BlueprintInfo struct {
	Title         string
	Source        *BlueprintRepoDetail
	Version       string                 `json:",omitempty" yaml:",omitempty"`
	ActuationTool BlueprintActuationTool `json:"actuationTool,omitempty" yaml:"actuationTool,omitempty"`
	Description   *BlueprintDescription  `json:",omitempty" yaml:",omitempty"`
	Icon          string                 `json:",omitempty" yaml:",omitempty"`
}

// BlueprintContent defines the detail for blueprint related content such as
// related documentation, diagrams, examples etc.
type BlueprintContent struct {
	// Diagrams are manually entered
	Diagrams      []BlueprintDiagram     `json:",omitempty" yaml:",omitempty"`
	Documentation []BlueprintListContent `json:",omitempty" yaml:",omitempty"`
	SubBlueprints []BlueprintMiscContent `json:"subBlueprints,omitempty" yaml:"subBlueprints,omitempty"`
	Examples      []BlueprintMiscContent `json:",omitempty" yaml:",omitempty"`
}

// BlueprintInterface the input and output variables for the blueprint
type BlueprintInterface struct {
	Variables []BlueprintVariable `json:",omitempty" yaml:",omitempty"`
	// VariableGroups are manually entered
	VariableGroups []BlueprintVariableGroup `json:"variableGroups,omitempty" yaml:"variableGroups,omitempty"`
	Outputs        []BlueprintOutput        `json:",omitempty" yaml:",omitempty"`
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
	Flavor  string `json:"type,omitempty" yaml:"type,omitempty"`
	Version string `json:",omitempty" yaml:",omitempty"`
}

type BlueprintDescription struct {
	Tagline      string   `json:",omitempty" yaml:",omitempty"`
	Detailed     string   `json:",omitempty" yaml:",omitempty"`
	PreDeploy    string   `json:"preDeploy,omitempty" yaml:"preDeploy,omitempty"`
	Architecture []string `json:"architecture,omitempty" yaml:"architecture,omitempty"`
}

type BlueprintMiscContent struct {
	Name     string
	Location string
}

// BlueprintDiagram is manually entered
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
	Description string      `json:",omitempty" yaml:",omitempty"`
	VarType     string      `yaml:"type"`
	Default     interface{} `json:",omitempty" yaml:",omitempty"`
	Required    bool
}

// BlueprintVariableGroup is manually entered
type BlueprintVariableGroup struct {
	Name        string
	Description string `json:",omitempty" yaml:",omitempty"`
	Variables   []string
}

type BlueprintOutput struct {
	Name        string
	Description string `json:",omitempty" yaml:",omitempty"`
}

type BlueprintRoles struct {
	Level string
	Roles []string
}
