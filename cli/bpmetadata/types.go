package bpmetadata

type BpMetadataDetail struct {
	Name           string
	Source         *BpRepoDetail
	Version        string
	ActuationTool  *BpActuationTool
	Title          string
	Description    *BpDescription
	SubBlueprints  []BpContent
	Examples       []BpContent
	Labels         []string
	Icon           string
	Diagrams       []BpDiagram
	Documentation  []BpDocumentation
	Variables      []BpVariable
	VariableGroups []BpVariableGroup
	Outputs        []BpOutputs
	Roles          []BpRoles
	Services       []string
}

type BpRepoDetail struct {
	Path       string
	SourceType string `yaml:"type"`
}

type BpActuationTool struct {
	Flavor  string `yaml:"type"`
	Version string
}

type BpDescription struct {
	Tagline   string
	Detailed  string
	PreDeploy string
}

type BpContent struct {
	Name     string
	Location string
}

type BpDiagram struct {
	Name        string
	AltText     string
	Description string
}

type BpDocumentation struct {
	Title string
	Url   string
}

type BpVariable struct {
	Name        string
	Description string
	VarType     string `yaml:"type"`
	Default     string
	Required    bool
}

type BpVariableGroup struct {
	Name        string
	Description string
	Variables   []string
}

type BpOutputs struct {
	Name        string
	Description string
}

type BpRoles struct {
	Granularity string
	Roles       []string
}
