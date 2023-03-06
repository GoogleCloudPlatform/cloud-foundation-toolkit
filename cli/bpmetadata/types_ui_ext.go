package bpmetadata

type ExtensionType string

const (
	Undefined ExtensionType = "UNDEFINED_TYPE"

	// General formats.
	EmailAddress    ExtensionType = "EMAIL_ADDRESS"
	MultiLineString ExtensionType = "MULTI_LINE_STRING"

	// GCE related.
	GceDiskImage       ExtensionType = "GCE_DISK_IMAGE"
	GceDiskType        ExtensionType = "GCE_DISK_TYPE"
	GceDiskSize        ExtensionType = "GCE_DISK_SIZE"
	GceMachineType     ExtensionType = "GCE_MACHINE_TYPE"
	GceNetwork         ExtensionType = "GCE_NETWORK"
	GceZone            ExtensionType = "GCE_ZONE"
	GceSubnetwork      ExtensionType = "GCE_SUBNETWORK"
	GceRegion          ExtensionType = "GCE_REGION"
	GceGpuType         ExtensionType = "GCE_GPU_TYPE"
	GceGpuCount        ExtensionType = "GCE_GPU_COUNT"
	GceExternalIp      ExtensionType = "GCE_EXTERNAL_IP"
	GceIpForwarding    ExtensionType = "GCE_IP_FORWARDING"
	GceFirewall        ExtensionType = "GCE_FIREWALL"
	GceFirewallRange   ExtensionType = "GCE_FIREWALL_RANGE"
	GceGenericResource ExtensionType = "GCE_GENERIC_RESOURCE"

	// GCS related.
	GcsBucket ExtensionType = "GCS_BUCKET"

	// IAM related.
	IamServiceAccount ExtensionType = "IAM_SERVICE_ACCOUNT"
)

// An extension for variables defined as part of DisplayVariable. The
// extension defines Google-specifc metadata necessary for choosing an
// appropriate input widget or adding restrictions to GCP-specific resources.
type GooglePropertyExtension struct {
	Type ExtensionType `yaml:"type"`

	// Some properties (e.g. GCE_MACHINE_TYPE) require a zone context in order to
	// determine the set of allowable values. This field references another
	// property from the schema, which must have type GCE_ZONE.
	ZoneProperty string `yaml:"zoneProperty,omitempty"`

	// Property-specific extensions
	GceMachineType    GceMachineTypeExtension     `yaml:"gceMachineType,omitempty"`
	GceDiskSize       GceDiskSizeExtension        `yaml:"gceDiskSize,omitempty"`
	GceSubnetwork     GceSubnetworkExtension      `yaml:"gceMachineType,omitempty"`
	GceResource       GceGenericResourceExtension `yaml:"gceSubnetwork,omitempty"`
	GceGpuType        GceGpuTypeExtension         `yaml:"gceGpuType,omitempty"`
	GceGpuCount       GceGpuCountExtension        `yaml:"gceGpuCount,omitempty"`
	GceNetwork        GceNetworkExtension         `yaml:"gceNetwork,omitempty"`
	GceExternalIp     GceExternalIpExtension      `yaml:"gceExternalIp,omitempty"`
	GceIpForwarding   GceIpForwardingExtension    `yaml:"gceIpForwarding,omitempty"`
	GceFirewall       GceFirewallExtension        `yaml:"gceFirewall,omitempty"`
	GceFirewallRange  GceFirewallRangeExtension   `yaml:"gceFirewallRange,omitempty"`
	GceZone           GceLocationExtension        `yaml:"gceZone,omitempty"`
	GceRegion         GceLocationExtension        `yaml:"gceRegion,omitempty"`
	IamServiceAccount IamServiceAccountExtension  `yaml:"iamServiceAccount,omitempty"`
}

type GceLocationExtension struct {
	WhitelistedZones   []string `yaml:"whitelistedZones,omitempty"`
	WhitelistedRegions []string `yaml:"whitelistedRegions,omitempty"`
}

type GceMachineTypeExtension struct {
	// Minimum cpu. Used to filter the list of selectable machine types.
	MinCpu int `yaml:"minCpu,omitempty"`

	// Minimum ram. Used to filter the list of selectable machine types.
	MinRamGb int `yaml:"minRamGb,omitempty"`

	// If true, custom machine types will not be selectable.
	// More info:
	// https://cloud.google.com/compute/docs/instances/creating-instance-with-custom-machine-type
	DisallowCustomMachineTypes bool `yaml:"disallowCustomMachineTypes,omitempty"`
}

type GceGpuTypeExtension struct {
	MachineType string `yaml:"machineType"`
	GpuType     string `yaml:"gpuType"`
}

type GceGpuCountExtension struct {
	// This field references another variable from the schema,
	// which must have type GceMachineType
	MachineTypeVariable string `yaml:"machineTypeVariable"`
}

type GceDiskSizeExtension struct {
	// The allowable range of disk sizes depends on the disk type. This field
	// references another variable from the schema, which must have type GceDiskType
	DiskTypeVariable string `yaml:"diskTypeVariable"`
}

type GceNetworkExtension struct {
	// AllowSharedVpcs indicates this solution can receive
	// shared VPC selflinks (fully qualified compute links).
	AllowSharedVpcs bool `yaml:"allowSharedVpcs"`
	// Used to indicate to which machine type this network interface will be
	// attached to.
	MachineTypeVariable string `yaml:"machineTypeVariable"`
}

type ExternalIpType string

const (
	IpEphemeral ExternalIpType = "EPHEMERAL"
	IpStaic     ExternalIpType = "STATIC"
)

type GceExternalIpExtension struct {
	// NetworkVariable is used to indicate the network variable in the schema
	// this external IP belongs to.
	NetworkVariable string `yaml:"networkVariable"`

	// Type specifies if the external IP is ephemeral or static.
	// Defaults to ephemeral if not specified.
	Type ExternalIpType `yaml:"externalIpType,omitempty"`
}

type GceIpForwardingExtension struct {
	// NetworkVariable is used to indicate the network variable in the schema
	// this external IP belongs to.
	NetworkVariable string `yaml:"networkVariable"`
	NotConfigurable bool   `yaml:"notConfigurable"`
}

type GceFirewallExtension struct {
	// NetworkVariable is used to indicate the network variable in the schema
	// this external IP belongs to.
	NetworkVariable string `yaml:"networkVariable"`
}

type GceFirewallRangeExtension struct {
	// FirewallVariable is used to indicate the firewall variable with the type
	// GceFirewall in the schema to which this firewall range belongs to.
	FirewallVariable string `yaml:"firewallVariable"`
}

type GceSubnetworkExtension struct {
	// Subnetwork variable requires a network context in order to determine the
	// set of available subnetworks. This field references another
	// variable from the schema, which must have type GceNetwork.
	NetworkVariable string `yaml:"networkVariable"`
}

type GceGenericResourceExtension struct {
	// GCE resource type to be fetched. This field references another
	// property from the schema, which must have type GceGenericResource.
	ResourceVariable string `yaml:"resourceVariable"`
}

type IamServiceAccountExtension struct {
	// List of IAM roles that to  grant to a new SA, or the roles to filter
	// existing SAs with
	Roles []string `yaml:"roles"`
}
