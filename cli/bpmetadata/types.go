package bpmetadata

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// BlueprintMetadata defines the overall structure for blueprint metadata.
// The cli command i.e. `cft blueprint metadata` attempts at auto-generating
// metadata if the blueprint is structured based on the TF blueprint template
// i.e. https://github.com/terraform-google-modules/terraform-google-module-template
// All fields within BlueprintMetadata and its children are denoted as:
// - Gen: auto-generated - <data source>
// - Gen: manually-authored
// - Gen: partial (contains child nodes that can be both auto-generated and manually authored)
type BlueprintMetadata struct {
	// Gen: auto-generated
	yaml.ResourceMeta `json:",inline" yaml:",inline"`
	// Gen: partial
	Spec BlueprintMetadataSpec `json:"spec" yaml:"spec"`
}

// BlueprintMetadataSpec defines the spec portion of the blueprint metadata.
type BlueprintMetadataSpec struct {
	// Gen: partial
	Info BlueprintInfo `json:"info,omitempty" yaml:"info,omitempty"`
	// Gen: partial
	Content BlueprintContent `json:"content,omitempty" yaml:"content,omitempty"`
	// Gen: partial
	Interfaces BlueprintInterface `json:"interfaces,omitempty" yaml:"interfaces,omitempty"`
	// Gen: auto-generated
	Requirements BlueprintRequirements `json:"requirements,omitempty" yaml:"requirements,omitempty"`
	// Gen: partial
	UI BlueprintUI `json:"ui,omitempty" yaml:"ui,omitempty"`
}

type BlueprintInfo struct {
	// Title for the blueprint.
	// Gen: auto-generated - First H1 text in readme.md.
	Title string `json:"title" yaml:"title"`

	// Blueprint source location and source type.
	// Gen: auto-generated - user will be prompted if repo information can not
	// be determined from the blueprint path.
	Source *BlueprintRepoDetail `json:"source,omitempty" yaml:"source,omitempty"`

	// Last released semantic version for the packaged blueprint.
	// Gen: auto-generated - From the `module_name` attribute of
	// the `provider_meta "google"` block.
	// E.g.
	// provider_meta "google" {
	//  module_name = "blueprints/terraform/terraform-google-log-analysis/v0.1.5"
	// }
	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	// Actuation tool e.g. Terraform and its required version.
	// Gen: auto-generated
	ActuationTool BlueprintActuationTool `json:"actuationTool,omitempty" yaml:"actuationTool,omitempty"`

	// Various types of descriptions associated with the blueprint.
	// Gen: auto-generated
	Description *BlueprintDescription `json:"description,omitempty" yaml:"description,omitempty"`

	// Path to an image representing the icon for the blueprint.
	// Will be set as "assets/icon.png", if present.
	// Gen: auto-generated
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty"`

	// The time estimate for configuring and deploying the blueprint.
	// Gen: auto-generated
	DeploymentDuration BlueprintTimeEstimate `json:"deploymentDuration,omitempty" yaml:"deploymentDuration,omitempty"`

	// The cost estimate for the blueprint based on preconfigured variables.
	// Gen: auto-generated
	CostEstimate BlueprintCostEstimate `json:"costEstimate,omitempty" yaml:"costEstimate,omitempty"`

	// A list of GCP cloud products used in the blueprint.
	// Gen: manually-authored
	CloudProducts []BlueprintCloudProduct `json:"cloudProducts,omitempty" yaml:"cloudProducts,omitempty"`

	// A list of GCP org policies to be checked for successful deployment.
	// Gen: manually-authored
	OrgPolicyChecks []BlueprintOrgPolicyCheck `json:"orgPolicyChecks,omitempty" yaml:"orgPolicyChecks,omitempty"`

	// A configuration of fixed and dynamic GCP quotas that apply to the blueprint.
	// Gen: manually-authored
	QuotaDetails []BlueprintQuotaDetail `json:"quotaDetails,omitempty" yaml:"quotaDetails,omitempty"`

	// Details on the author producing the blueprint.
	// Gen: manually-authored
	Author BlueprintAuthor `json:"author,omitempty" yaml:"author,omitempty"`

	// Details on software installed as part of the blueprint.
	// Gen: manually-authored
	SoftwareGroups []BlueprintSoftwareGroup `json:"softwareGroups,omitempty" yaml:"softwareGroups,omitempty"`

	// Support offered, if any for the blueprint.
	// Gen: manually-authored
	SupportInfo BlueprintSupport `json:"supportInfo,omitempty" yaml:"supportInfo,omitempty"`

	// Specifies if the blueprint supports single or multiple deployments per GCP project.
	// If set to true, the blueprint can not be deployed more than once in the same GCP project.
	// Gen: manually-authored
	SingleDeployment bool `json:"singleDeployment,omitempty" yaml:"singleDeployment,omitempty"`
}

type BlueprintRepoDetail struct {
	// Gen: auto-generated - URL from the .git dir.
	// Can be manually overridden with a custom URL if needed.
	Repo string `json:"repo" yaml:"repo"`

	// Gen: auto-generated - set as "git" for now until more
	// types are supported.
	SourceType string `json:"sourceType" yaml:"sourceType"`

	// Gen: auto-generated - not set for root modules but
	// set as the module name for submodules, if found.
	Dir string `json:"dir,omitempty" yaml:"dir,omitempty"`
}

type BlueprintActuationTool struct {
	// Gen: auto-generated - set as "Terraform" for now until
	//more flavors are supported.
	Flavor string `json:"flavor,omitempty" yaml:"flavor,omitempty"`

	// Required version for the actuation tool.
	// Gen: auto-generated - For Terraform this is the `required_version`
	// set in `terraform` block. E.g.
	// terraform {
	//   required_version = ">= 0.13"
	// }
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// All descriptions are set with the markdown content immediately
// after each type's heading declaration in readme.md.
type BlueprintDescription struct {
	// Gen: auto-generated - Markdown after "### Tagline".
	Tagline string `json:"tagline,omitempty" yaml:"tagline,omitempty"`

	// Gen: auto-generated - Markdown after "### Detailed".
	Detailed string `json:"detailed,omitempty" yaml:"detailed,omitempty"`

	// Gen: auto-generated - Markdown after "### PreDeploy".
	PreDeploy string `json:"preDeploy,omitempty" yaml:"preDeploy,omitempty"`

	// Gen: auto-generated - Markdown after "### Html".
	HTML string `json:"html,omitempty" yaml:"html,omitempty"`

	// Gen: auto-generated - Markdown after "### EulaUrls".
	EulaURLs []string `json:"eulaUrls,omitempty" yaml:"eulaUrls,omitempty"`

	// Gen: auto-generated - Markdown after "### Architecture"
	// Deprecated. Use BlueprintContent.Architecture instead.
	Architecture []string `json:"architecture,omitempty" yaml:"architecture,omitempty"`
}

// A time estimate in secs required for configuring and deploying the blueprint.
type BlueprintTimeEstimate struct {
	// Gen: auto-generated - Set using the content defined under "### DeploymentTime" E.g.
	// ### DeploymentTime
	// - Configuration: X secs
	// - Deployment: Y secs
	ConfigurationSecs int `json:"configurationSecs,omitempty" yaml:"configurationSecs,omitempty"`
	DeploymentSecs    int `json:"deploymentSecs,omitempty" yaml:"deploymentSecs,omitempty"`
}

// The cost estimate for the blueprint based on pre-configured variables.
type BlueprintCostEstimate struct {
	// Gen: auto-generated - Set using the content defined under "### Cost" as a link
	// with a description E.g.
	// ### Cost
	// [View cost details](https://cloud.google.com/products/calculator?hl=en_US&_ga=2.1665458.-226505189.1675191136#id=02fb0c45-cc29-4567-8cc6-f72ac9024add)
	Description string `json:"description" yaml:"description"`
	URL         string `json:"url" yaml:"url"`
}

// GCP cloud product(s) used in the blueprint.
type BlueprintCloudProduct struct {
	// A top-level (e.g. "Compute Engine") or secondary (e.g. "Binary Authorization")
	// product used in the blueprint.
	// Gen: manually-authored
	ProductId string `json:"productId,omitempty" yaml:"productId,omitempty"`

	// Url for the product.
	// Gen: manually-authored
	PageURL string `json:"pageUrl" yaml:"pageUrl"`

	// A label string for the product, if it is not an integrated GCP product.
	// E.g. "Data Studio"
	// Gen: manually-authored
	Label string `json:"label,omitempty" yaml:"label,omitempty"`

	// Is the product's landing page external to the GCP console e.g.
	// lookerstudio.google.com
	// Gen: manually-authored
	IsExternal bool `json:"isExternal,omitempty" yaml:"isExternal,omitempty"`
}

// BlueprintOrgPolicyCheck defines GCP org policies to be checked
// for successful deployment
type BlueprintOrgPolicyCheck struct {
	// Id for the policy e.g. "compute-vmExternalIpAccess"
	// Gen: manually-authored
	PolicyId string `json:"policyId" yaml:"policyId"`

	// If not set, it is assumed any version of this org policy
	// prevents successful deployment of this solution.
	// Gen: manually-authored
	RequiredValues []string `json:"requiredValues,omitempty" yaml:"requiredValues,omitempty"`
}

type QuotaResourceType string

const (
	QuotaResTypeUndefined   QuotaResourceType = "QRT_UNDEFINED"
	QuotaResTypeGCEInstance QuotaResourceType = "QRT_GCE_INSTANCE"
	QuotaResTypeGCEDisk     QuotaResourceType = "QRT_GCE_DISK"
)

type QuotaType string

const (
	MachineType QuotaType = "MACHINE_TYPE"
	CPUs        QuotaType = "CPUs"
	DiskType    QuotaType = "DISK_TYPE"
	DiskSizeGB  QuotaType = "SIZE_GB"
)

type BlueprintQuotaDetail struct {
	// DynamicVariable, if provided, associates the provided input variable
	// with the corresponding resource and quota type. In its absence, the quota
	// detail is assumed to be fixed.
	// Gen: manually-authored
	DynamicVariable string `json:"dynamicVariable,omitempty" yaml:"dynamicVariable,omitempty"`

	// ResourceType is the type of resource the quota will be applied to i.e.
	// GCE Instance or Disk etc.
	// Gen: manually-authored
	ResourceType QuotaResourceType `json:"resourceType" yaml:"resourceType" jsonschema:"enum=QRT_GCE_INSTANCE,enum=QRT_GCE_DISK,enum=QRT_UNDEFINED"`

	// QuotaType is a key/value pair of the actual quotas and their corresponding
	// values.
	// Gen: manually-authored
	QuotaType map[QuotaType]string `json:"quotaType" yaml:"quotaType"`
}

type BlueprintAuthor struct {
	// Name of template author or organization.
	// Gen: manually-authored
	Title string `json:"title" yaml:"title"`

	// Description of the author.
	// Gen: manually-authored
	Description string `json:"description,omitempty"  yaml:"description,omitempty"`

	// Link to the author's website.
	// Gen: manually-authored
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
}

type SoftwareGroupType string

const (
	SG_Unspecified SoftwareGroupType = "SG_UNSPECIFIED"
	SG_OS          SoftwareGroupType = "SG_OS"
)

// A group of related software components for the blueprint.
type BlueprintSoftwareGroup struct {
	// Pre-defined software types.
	// Gen: manually-authored
	Type SoftwareGroupType `json:"type,omitempty" yaml:"type,omitempty" jsonschema:"enum=SG_UNSPECIFIED,enum=SG_OS"`

	// Software components belonging to this group.
	// Gen: manually-authored
	Software []BlueprintSoftware `json:"software,omitempty" yaml:"software,omitempty"`
}

// A description of a piece of a single software component
// installed by the blueprint.
type BlueprintSoftware struct {
	// User-visible title.
	// Gen: manually-authored
	Title string `json:"title" yaml:"title"`

	// Software version.
	// Gen: manually-authored
	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	// Link to development site or marketing page for this software.
	// Gen: manually-authored
	URL string `json:"url,omitempty" yaml:"url,omitempty"`

	// Link to license page.
	// Gen: manually-authored
	LicenseURL string `json:"licenseUrl,omitempty" yaml:"licenseUrl,omitempty"`
}

// A description of a support option
type BlueprintSupport struct {
	// Description of the support option.
	// Gen: manually-authored
	Description string `json:"description" yaml:"description"`

	// Link to the page providing this support option.
	// Gen: manually-authored
	URL string `json:"url,omitempty" yaml:"url,omitempty"`

	// The organization or group that provides the support option (e.g.:
	// "Community", "Google").
	// Gen: manually-authored
	Entity string `json:"entity,omitempty" yaml:"entity,omitempty"`

	// Whether to show the customer's support ID.
	// Gen: manually-authored
	ShowSupportId bool `json:"showSupportId,omitempty" yaml:"showSupportId,omitempty"`
}

// BlueprintContent defines the detail for blueprint related content such as
// related documentation, diagrams, examples etc.
type BlueprintContent struct {
	// Gen: auto-generated
	Architecture BlueprintArchitecture `json:"architecture,omitempty" yaml:"architecture,omitempty"`

	// Gen: manually-authored
	Diagrams []BlueprintDiagram `json:"diagrams,omitempty" yaml:"diagrams,omitempty"`

	// Gen: auto-generated - the list content following the "## Documentation" tag. E.g.
	// ## Documentation
	// - [Hosting a Static Website](https://cloud.google.com/storage/docs/hosting-static-website)
	Documentation []BlueprintListContent `json:"documentation,omitempty" yaml:"documentation,omitempty"`

	// Gen: auto-generated - blueprints under the modules/ folder.
	SubBlueprints []BlueprintMiscContent `json:"subBlueprints,omitempty" yaml:"subBlueprints,omitempty"`

	// Gen: auto-generated - examples under the examples/ folder.
	Examples []BlueprintMiscContent `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// BlueprintInterface defines the input and output variables for the blueprint.
type BlueprintInterface struct {
	// Gen: auto-generated - all defined variables for the blueprint
	Variables []BlueprintVariable `json:"variables,omitempty" yaml:"variables,omitempty"`

	// Gen: manually-authored
	VariableGroups []BlueprintVariableGroup `json:"variableGroups,omitempty" yaml:"variableGroups,omitempty"`

	// Gen: auto-generated - all defined outputs for the blueprint
	Outputs []BlueprintOutput `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

// BlueprintRequirements defines the roles required and the associated services
// that need to be enabled to provision blueprint resources.
type BlueprintRequirements struct {
	// Gen: auto-generated - all roles required for the blueprint in test/setup/iam.tf
	// as the "int_required_roles" local. E.g.
	// locals {
	//   int_required_roles = [
	//     "roles/compute.admin",
	//   ]
	// }
	Roles []BlueprintRoles `json:"roles,omitempty" yaml:"roles,omitempty"`

	// Gen: auto-generated - all services required for the blueprint in test/setup/main.tf
	// as "activate_apis" in the project module.
	Services []string `json:"services,omitempty" yaml:"services,omitempty"`
}

type BlueprintArchitecture struct {
	// Gen: auto-generated - the URL & list content following the "## Architecture" tag e.g.
	// ## Architecture
	// ![Blueprint Architecture](assets/architecture.png)
	// 1. Step no. 1
	// 2. Step no. 2
	// 3. Step no. 3
	DiagramURL string `json:"diagramUrl" yaml:"diagramUrl"`

	// Gen: auto-generated - the list items following the "## Architecture" tag.
	Description []string `json:"description" yaml:"description"`
}

type BlueprintDiagram struct {
	Name        string `json:"name" yaml:"name"`
	AltText     string `json:"altText,omitempty" yaml:"altText,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type BlueprintMiscContent struct {
	Name     string `json:"name" yaml:"name"`
	Location string `json:"location,omitempty" yaml:"location,omitempty"`
}

type BlueprintListContent struct {
	Title string `json:"title" yaml:"title"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
}

type BlueprintVariable struct {
	Name         string      `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string      `json:"description,omitempty" yaml:"description,omitempty"`
	VarType      string      `json:"varType,omitempty" yaml:"varType,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	Required     bool        `json:"required,omitempty" yaml:"required,omitempty"`
}

type BlueprintVariableGroup struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   []string `json:"variables,omitempty" yaml:"variables,omitempty"`
}

type BlueprintOutput struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type BlueprintRoles struct {
	Level string   `json:"level" yaml:"level"`
	Roles []string `json:"roles" yaml:"roles"`
}
