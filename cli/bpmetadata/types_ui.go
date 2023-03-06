package bpmetadata

// BlueprintUI is the top-level structure for holding UI specific metadata
type BlueprintUI struct {
	// The top-level input section that defines the list of properties and
	// their sections on the deployment page
	Input BlueprintUIInput `yaml:"input"`

	// The top-level section for listing runtime (or solution output) information
	// i.e. the console URL for the VM or a button to ssh into the VM etc based on
	Runtime BlueprintUIOutput `yaml:"runtime"`

	// The solution version corresponding to the version in the
	// standard metadata file for keeping data in sync.
	Version string `yaml:"version"`
}

// BlueprintUIInput is the structure for holding Input and Input Section (i.e. groups) specific metadata
type BlueprintUIInput struct {

	// Properties is a map defining all inputs on the UI
	DisplayVariables map[string]DisplayVariable `yaml:"variables"`

	// Sections is a generic structure for grouping inputs together
	DisplaySections []DisplaySection `yaml:"sections"`
}

// Additional display specific metadata pertaining to a particular
// input variable.
type DisplayVariable struct {
	// The variable name from the corresponding standard metadata file.
	Name string `yaml:"name"`

	// Visible title for the variable on the UI.
	Title bool `yaml:"title,omitempty"`

	// A flag to hide or show the property on the UI.
	Visible bool `yaml:"visible,omitempty"`

	// Property tooltip.
	Tooltip string `yaml:"tooltip,omitempty"`

	// Placeholder text (when there is no default).
	Placeholder string `yaml:"placeholder,omitempty"`

	// Text describing the validation rules for the property. Typically shown
	// after an invalid input.
	Validation string `yaml:"validation,omitempty"`

	// The pattern for the value for the input variable
	Pattern string `yaml:"pattern,omitempty"`

	// Minimum no. of values for the input variable
	Minimum int `yaml:"min,omitempty"`

	// Max no. of values for the input variable
	Maximum int `yaml:"max,omitempty"`

	// The name of a section to which this property belongs.
	// Properties belong to the root section if this field is
	// not set.
	Section string `yaml:"section,omitempty"`

	// Designates that this property has no impact on the costs, quotas, or
	// permissions associated with the resources in the expanded deployment.
	// Typically true for application-specific properties that do not affect the
	// size or number of instances in the deployment.
	NoResourceImpact bool `yaml:"noResourceImpact,omitempty"`

	// UI property extension associated with the input variable
	UIPropertyExtension GooglePropertyExtension `yaml:"x-googleProperty,omitempty"`
}

// A logical group of properties. [Section][]s may also be grouped into
// sub-sections. Child of [Input][]
type DisplaySection struct {
	// The name of the section, referenced by DisplayVariable.Section
	// Section names must be unique.
	Name string `yaml:"name"`

	// Section title.
	// If not provided, name will be used instead
	Title string `yaml:"title,omitempty"`

	// Section tooltip.
	Tooltip string `yaml:"tooltip,omitempty"`

	// Section subtext.
	Subtext string `yaml:"subtext,omitempty"`

	// The name of the parent section (if parent is not the root section).
	Parent string `yaml:"parent,omitempty"`
}

type BlueprintUIOutput struct {
	// Short message to be displayed while the solution is deploying.
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
	// Boolean expression.
	ShowIf string `yaml:"showIf,omitempty"`
}
