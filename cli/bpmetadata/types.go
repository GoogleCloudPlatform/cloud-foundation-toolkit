package bpmetadata

import "sigs.k8s.io/kustomize/kyaml/yaml"

// BlueprintMetadata defines the overall structure for blueprint metadata details
type BlueprintMetadata struct {
	Meta yaml.ResourceMeta
	Spec *BlueprintMetadataSpec `yaml:"spec"`
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
	Title         string               `yaml:"title"`
	Source        *BlueprintRepoDetail `yaml:"source"`
	Version       string
	ActuationTool *BlueprintActuationTool
	Description   *BlueprintDescription
	Icon          string
}

// BlueprintContent defines the detail for blueprint related content such as
// related documentation, diagrams, examples etc.
type BlueprintContent struct {
	Diagrams      []BlueprintDiagram
	Documentation []BlueprintListContent
	SubBlueprints []BlueprintMiscContent
	Examples      []BlueprintMiscContent
}

// BlueprintInterface the input and output variables for the blueprint
type BlueprintInterface struct {
	Variables      []BlueprintVariable
	VariableGroups []BlueprintVariableGroup
	Outputs        []BlueprintOutput
}

// BlueprintRequirements defines the roles required and the assocaited services
// that need to be enabled to provision blueprint resources
type BlueprintRequirements struct {
	Roles    []BlueprintRoles
	Services []string
}

type BlueprintRepoDetail struct {
	Repo       string `yaml:"repo"`
	SourceType string `yaml:"sourceType"`
}

type BlueprintActuationTool struct {
	Flavor  string `yaml:"type"`
	Version string
}

type BlueprintDescription struct {
	Tagline   string
	Detailed  string
	PreDeploy string
}

type BlueprintMiscContent struct {
	Name     string
	Location string
}

type BlueprintDiagram struct {
	Name        string
	AltText     string
	Description string
}

type BlueprintListContent struct {
	Title string
	Url   string
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
	Granularity string
	Roles       []string
}
