package bpmetadata

// BlueprintUI is the top-level structure for holding UI specific metadata.
type BlueprintUI struct {
	// The top-level input section that defines the list of variables and
	// their sections on the deployment page.
	Input BlueprintUIInput `yaml:"input"`

	// The top-level section for listing runtime (or blueprint output) information
	// i.e. the console URL for the VM or a button to ssh into the VM etc based on.
	Runtime BlueprintUIOutput `yaml:"runtime"`
}

// BlueprintUIInput is the structure for holding Input and Input Section (i.e. groups) specific metadata.
type BlueprintUIInput struct {

	// variables is a map defining all inputs on the UI.
	DisplayVariables map[string]DisplayVariable `yaml:"variables"`

	// Sections is a generic structure for grouping inputs together.
	DisplaySections []DisplaySection `yaml:"sections"`
}

// Additional display specific metadata pertaining to a particular
// input variable.
type DisplayVariable struct {
	// The variable name from the corresponding standard metadata file.
	Name string `yaml:"name"`

	// Visible title for the variable on the UI.
	Title bool `yaml:"title,omitempty"`

	// A flag to hide or show the variable on the UI.
	Visible bool `yaml:"visible,omitempty"`

	// Variable tooltip.
	Tooltip string `yaml:"tooltip,omitempty"`

	// Placeholder text (when there is no default).
	Placeholder string `yaml:"placeholder,omitempty"`

	// Text describing the validation rules for the variable based
	// on a regular expression.
	// Typically shown after an invalid input.
	RegExValidation string `yaml:"regexValidation,omitempty"`

	// Minimum no. of values for the input variable.
	Minimum int `yaml:"min,omitempty"`

	// Max no. of values for the input variable.
	Maximum int `yaml:"max,omitempty"`

	// The name of a section to which this variable belongs.
	// variables belong to the root section if this field is
	// not set.
	Section string `yaml:"section,omitempty"`

	// Designates that this variable has no impact on the costs, quotas, or
	// permissions associated with the resources in the expanded deployment.
	// Typically true for application-specific variables that do not affect the
	// size or number of instances in the deployment.
	ResourceImpact bool `yaml:"resourceImpact,omitempty"`

	// UI extension associated with the input variable.
	// E.g. for rendering a GCE machine type selector:
	//
	// x-googleProperty:
	//   type: GCE_MACHINE_TYPE
	//   zoneProperty: myZone
	//   gceMachineType:
	//     minCpu: 2
	//     minRamGb: 6
	UIDisplayVariableExtension GooglePropertyExtension `yaml:"x-googleProperty,omitempty"`
}

// A logical group of variables. [Section][]s may also be grouped into
// sub-sections.
type DisplaySection struct {
	// The name of the section, referenced by DisplayVariable.Section
	// Section names must be unique.
	Name string `yaml:"name"`

	// Section title.
	// If not provided, name will be used instead.
	Title string `yaml:"title,omitempty"`

	// Section tooltip.
	Tooltip string `yaml:"tooltip,omitempty"`

	// Section subtext.
	Subtext string `yaml:"subtext,omitempty"`

	// The name of the parent section (if parent is not the root section).
	Parent string `yaml:"parent,omitempty"`
}

type BlueprintUIOutput struct {
	// Short message to be displayed while the blueprint is deploying.
	// At most 128 characters.
	OutputMessage string `yaml:"outputMessage,omitempty"`

	// List of suggested actions to take.
	SuggestedActions []UIActionItem `yaml:"suggestedActions,omitempty"`
}

// An item appearing in a list of required or suggested steps.
type UIActionItem struct {
	// Summary heading for the item.
	// Required. Accepts string expressions. At most 64 characters.
	Heading string `yaml:"heading"`

	// Longer description of the item.
	// At least one description or snippet is required.
	// Accepts string expressions. HTML <code>&lt;a href&gt;</code>
	// tags only. At most 512 characters.
	Description string `yaml:"description"`

	// Fixed-width formatted code snippet.
	// At least one description or snippet is required.
	// Accepts string expressions. UTF-8 text. At most 512 characters.
	Snippet string `yaml:"snippet"`

	// If present, this expression determines whether the item is shown.
	// Should be in the form of a Boolean expression e.g. outputs().hasExternalIP
	// where `externalIP` is the output.
	ShowIf string `yaml:"showIf,omitempty"`
}
