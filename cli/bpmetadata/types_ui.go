package bpmetadata

// BlueprintUI is the top-level structure for holding UI specific metadata.
type BlueprintUI struct {
	// The top-level input section that defines the list of variables and
	// their sections on the deployment page.
	Input BlueprintUIInput `json:"input,omitempty" yaml:"input,omitempty"`

	// The top-level section for listing runtime (or blueprint output) information
	// i.e. the console URL for the VM or a button to ssh into the VM etc based on.
	Runtime BlueprintUIOutput `json:"runtime,omitempty" yaml:"runtime,omitempty"`
}

// BlueprintUIInput is the structure for holding Input and Input Section (i.e. groups) specific metadata.
type BlueprintUIInput struct {

	// variables is a map defining all inputs on the UI.
	DisplayVariables map[string]*DisplayVariable `json:"variables,omitempty" yaml:"variables,omitempty"`

	// Sections is a generic structure for grouping inputs together.
	DisplaySections []DisplaySection `json:"sections,omitempty" yaml:"sections,omitempty"`
}

// Additional display specific metadata pertaining to a particular
// input variable.
type DisplayVariable struct {
	// The variable name from the corresponding standard metadata file.
	Name string `json:"name" yaml:"name"`

	// Visible title for the variable on the UI. If not present,
	// Name will be used for the Title.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// A flag to hide or show the variable on the UI.
	Invisible bool `json:"invisible,omitempty" yaml:"invisible,omitempty"`

	// Variable tooltip.
	Tooltip string `json:"tooltip,omitempty" yaml:"tooltip,omitempty"`

	// Placeholder text (when there is no default).
	Placeholder string `json:"placeholder,omitempty" yaml:"placeholder,omitempty"`

	// Text describing the validation rules for the variable based
	// on a regular expression.
	// Typically shown after an invalid input.
	RegExValidation string `json:"regexValidation,omitempty" yaml:"regexValidation,omitempty"`

	// Minimum no. of inputs for the input variable.
	MinimumItems int `json:"minItems,omitempty" yaml:"minItems,omitempty"`

	// Max no. of inputs for the input variable.
	MaximumItems int `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`

	// Minimum length for string values.
	MinimumLength int `json:"minLength,omitempty" yaml:"minLength,omitempty"`

	// Max length for string values.
	MaximumLength int `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`

	// Minimum value for numeric types.
	Minimum int `json:"min,omitempty" yaml:"min,omitempty"`

	// Max value for numeric types.
	Maximum int `json:"max,omitempty" yaml:"max,omitempty"`

	// The name of a section to which this variable belongs.
	// variables belong to the root section if this field is
	// not set.
	Section string `json:"section,omitempty" yaml:"section,omitempty"`

	// UI extension associated with the input variable.
	// E.g. for rendering a GCE machine type selector:
	//
	// x-googleProperty:
	//   type: GCE_MACHINE_TYPE
	//   zoneProperty: myZone
	//   gceMachineType:
	//     minCpu: 2
	//     minRamGb: 6
	UIDisplayVariableExtension GooglePropertyExtension `json:"x-googleProperty,omitempty" yaml:"x-googleProperty,omitempty"`
}

// A logical group of variables. [Section][]s may also be grouped into
// sub-sections.
type DisplaySection struct {
	// The name of the section, referenced by DisplayVariable.Section
	// Section names must be unique.
	Name string `json:"name" yaml:"name"`

	// Section title.
	// If not provided, name will be used instead.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// Section tooltip.
	Tooltip string `json:"tooltip,omitempty" yaml:"tooltip,omitempty"`

	// Section subtext.
	Subtext string `json:"subtext,omitempty" yaml:"subtext,omitempty"`

	// The name of the parent section (if parent is not the root section).
	Parent string `json:"parent,omitempty" yaml:"parent,omitempty"`
}

type BlueprintUIOutput struct {
	// Short message to be displayed while the blueprint is deploying.
	// At most 128 characters.
	OutputMessage string `json:"outputMessage,omitempty" yaml:"outputMessage,omitempty"`

	// List of suggested actions to take.
	SuggestedActions []UIActionItem `json:"suggestedActions,omitempty" yaml:"suggestedActions,omitempty"`
}

// An item appearing in a list of required or suggested steps.
type UIActionItem struct {
	// Summary heading for the item.
	// Required. Accepts string expressions. At most 64 characters.
	Heading string `json:"heading" yaml:"heading"`

	// Longer description of the item.
	// At least one description or snippet is required.
	// Accepts string expressions. HTML <code>&lt;a href&gt;</code>
	// tags only. At most 512 characters.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Fixed-width formatted code snippet.
	// At least one description or snippet is required.
	// Accepts string expressions. UTF-8 text. At most 512 characters.
	Snippet string `json:"snippet,omitempty" yaml:"snippet,omitempty"`

	// If present, this expression determines whether the item is shown.
	// Should be in the form of a Boolean expression e.g. outputs().hasExternalIP
	// where `externalIP` is the output.
	ShowIf string `json:"showIf,omitempty" yaml:"showIf,omitempty"`
}
