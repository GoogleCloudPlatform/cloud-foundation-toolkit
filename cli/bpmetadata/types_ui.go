package bpmetadata

// BlueprintUI is the top-level structure for holding UI specific metadata.
type BlueprintUI struct {
	// The top-level input section that defines the list of variables and
	// their sections on the deployment page.
	// Gen: partial
	Input BlueprintUIInput `json:"input,omitempty" yaml:"input,omitempty"`

	// The top-level section for listing runtime (or blueprint output) information
	// i.e. the console URL for the VM or a button to ssh into the VM etc based on.
	// Gen: manually-authored
	Runtime BlueprintUIOutput `json:"runtime,omitempty" yaml:"runtime,omitempty"`
}

// BlueprintUIInput is the structure for holding Input and Input Section (i.e. groups) specific metadata.
type BlueprintUIInput struct {
	// variables is a map defining all inputs on the UI.
	// Gen: partial
	Variables map[string]*DisplayVariable `json:"variables,omitempty" yaml:"variables,omitempty"`

	// Sections is a generic structure for grouping inputs together.
	// Gen: manually-authored
	Sections []DisplaySection `json:"sections,omitempty" yaml:"sections,omitempty"`
}

// Additional display specific metadata pertaining to a particular
// input variable.
type DisplayVariable struct {
	// The variable name from the corresponding standard metadata file.
	// Gen: auto-generated - the Terraform variable name
	Name string `json:"name" yaml:"name"`

	// Visible title for the variable on the UI. If not present,
	// Name will be used for the Title.
	// Gen: auto-generated - the Terraform variable converted to title case e.g.
	// variable "bucket_admins" will convert to "Bucket Admins" as the title.
	Title string `json:"title" yaml:"title"`

	// A flag to hide or show the variable on the UI.
	// Gen: manually-authored
	Invisible bool `json:"invisible,omitempty" yaml:"invisible,omitempty"`

	// Variable tooltip.
	// Gen: manually-authored
	Tooltip string `json:"tooltip,omitempty" yaml:"tooltip,omitempty"`

	// Placeholder text (when there is no default).
	// Gen: manually-authored
	Placeholder string `json:"placeholder,omitempty" yaml:"placeholder,omitempty"`

	// Text describing the validation rules for the property. Typically shown
	// after an invalid input.
	// Optional. UTF-8 text. No markup. At most 128 characters.
	// Gen: manually-authored
	Validation string `json:"validation,omitempty" yaml:"validation,omitempty"`

	// Regex based validation rules for the variable.
	// Gen: manually-authored
	RegExValidation string `json:"regexValidation,omitempty" yaml:"regexValidation,omitempty"`

	// Minimum no. of inputs for the input variable.
	// Gen: manually-authored
	MinimumItems int `json:"minItems,omitempty" yaml:"minItems,omitempty"`

	// Max no. of inputs for the input variable.
	// Gen: manually-authored
	MaximumItems int `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`

	// Minimum length for string values.
	// Gen: manually-authored
	MinimumLength int `json:"minLength,omitempty" yaml:"minLength,omitempty"`

	// Max length for string values.
	// Gen: manually-authored
	MaximumLength int `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`

	// Minimum value for numeric types.
	// Gen: manually-authored
	Minimum float32 `json:"min,omitempty" yaml:"min,omitempty"`

	// Max value for numeric types.
	// Gen: manually-authored
	Maximum float32 `json:"max,omitempty" yaml:"max,omitempty"`

	// The name of a section to which this variable belongs.
	// variables belong to the root section if this field is
	// not set.
	// Gen: manually-authored
	Section string `json:"section,omitempty" yaml:"section,omitempty"`

	// UI extension associated with the input variable.
	// E.g. for rendering a GCE machine type selector:
	//
	// xGoogleProperty:
	//   type: GCE_MACHINE_TYPE
	//   zoneProperty: myZone
	//   gceMachineType:
	//     minCpu: 2
	//     minRamGb: 6
	// Gen: manually-authored
	XGoogleProperty GooglePropertyExtension `json:"xGoogleProperty,omitempty" yaml:"xGoogleProperty,omitempty"`
}

// A logical group of variables. [Section][]s may also be grouped into
// sub-sections.
type DisplaySection struct {
	// The name of the section, referenced by DisplayVariable.Section
	// Section names must be unique.
	// Gen: manually-authored
	Name string `json:"name" yaml:"name"`

	// Section title.
	// If not provided, name will be used instead.
	// Gen: manually-authored
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// Section tooltip.
	// Gen: manually-authored
	Tooltip string `json:"tooltip,omitempty" yaml:"tooltip,omitempty"`

	// Section subtext.
	// Gen: manually-authored
	Subtext string `json:"subtext,omitempty" yaml:"subtext,omitempty"`

	// The name of the parent section (if parent is not the root section).
	// Gen: manually-authored
	Parent string `json:"parent,omitempty" yaml:"parent,omitempty"`
}

type BlueprintUIOutput struct {
	// Short message to be displayed while the blueprint is deploying.
	// At most 128 characters.
	// Gen: manually-authored
	OutputMessage string `json:"outputMessage,omitempty" yaml:"outputMessage,omitempty"`

	// List of suggested actions to take.
	// Gen: manually-authored
	SuggestedActions []UIActionItem `json:"suggestedActions,omitempty" yaml:"suggestedActions,omitempty"`

	// Outputs is a map defining a subset of Terraform outputs on the UI
	// that may need additional UI configuration.
	// Gen: manually-authored
	Outputs map[string]DisplayOutput `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

// An item appearing in a list of required or suggested steps.
type UIActionItem struct {
	// Summary heading for the item.
	// Required. Accepts string expressions. At most 64 characters.
	// Gen: manually-authored
	Heading string `json:"heading" yaml:"heading"`

	// Longer description of the item.
	// At least one description or snippet is required.
	// Accepts string expressions. HTML <code>&lt;a href&gt;</code>
	// tags only. At most 512 characters.
	// Gen: manually-authored
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Fixed-width formatted code snippet.
	// At least one description or snippet is required.
	// Accepts string expressions. UTF-8 text. At most 512 characters.
	// Gen: manually-authored
	Snippet string `json:"snippet,omitempty" yaml:"snippet,omitempty"`

	// If present, this expression determines whether the item is shown.
	// Should be in the form of a Boolean expression e.g. outputs().hasExternalIP
	// where `externalIP` is the output.
	// Gen: manually-authored
	ShowIf string `json:"showIf,omitempty" yaml:"showIf,omitempty"`
}

// Additional display specific metadata pertaining to a particular
// Terraform output.
type DisplayOutput struct {
	// OpenInNewTab defines if the Output action should be opened
	// in a new tab.
	// Gen: manually-authored
	OpenInNewTab bool `json:"openInNewTab,omitempty" yaml:"openInNewTab,omitempty"`

	// ShowInNotification defines if the Output should shown in
	// notification for the deployment.
	// Gen: manually-authored
	ShowInNotification bool `json:"showInNotification,omitempty" yaml:"showInNotification,omitempty"`
}
