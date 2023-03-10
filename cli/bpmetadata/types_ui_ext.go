package bpmetadata

type ExtensionType string

const (
	Undefined ExtensionType = "UNDEFINED_TYPE"

	// General formats.
	EmailAddress    ExtensionType = "EMAIL_ADDRESS"
	MultiLineString ExtensionType = "MULTI_LINE_STRING"

	// GCE related.
	GCEDiskImage       ExtensionType = "GCE_DISK_IMAGE"
	GCEDiskType        ExtensionType = "GCE_DISK_TYPE"
	GCEDiskSize        ExtensionType = "GCE_DISK_SIZE"
	GCEMachineType     ExtensionType = "GCE_MACHINE_TYPE"
	GCENetwork         ExtensionType = "GCE_NETWORK"
	GCEZone            ExtensionType = "GCE_ZONE"
	GCESubnetwork      ExtensionType = "GCE_SUBNETWORK"
	GCERegion          ExtensionType = "GCE_REGION"
	GCEGPUType         ExtensionType = "GCE_GPU_TYPE"
	GCEGPUCount        ExtensionType = "GCE_GPU_COUNT"
	GCEExternalIP      ExtensionType = "GCE_EXTERNAL_IP"
	GCEIPForwarding    ExtensionType = "GCE_IP_FORWARDING"
	GCEFirewall        ExtensionType = "GCE_FIREWALL"
	GCEFirewallRange   ExtensionType = "GCE_FIREWALL_RANGE"
	GCEGenericResource ExtensionType = "GCE_GENERIC_RESOURCE"

	// GCS related.
	GCSBucket ExtensionType = "GCS_BUCKET"

	// IAM related.
	IAMServiceAccount ExtensionType = "IAM_SERVICE_ACCOUNT"
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

	// Property-specific extensions.
	GCEMachineType    GCEMachineTypeExtension     `yaml:"gceMachineType,omitempty"`
	GCEDiskSize       GCEDiskSizeExtension        `yaml:"gceDiskSize,omitempty"`
	GCESubnetwork     GCESubnetworkExtension      `yaml:"gceMachineType,omitempty"`
	GCEResource       GCEGenericResourceExtension `yaml:"gceSubnetwork,omitempty"`
	GCEGPUType        GCEGPUTypeExtension         `yaml:"gceGpuType,omitempty"`
	GCEGPUCount       GCEGPUCountExtension        `yaml:"gceGpuCount,omitempty"`
	GCENetwork        GCENetworkExtension         `yaml:"gceNetwork,omitempty"`
	GCEExternalIP     GCEExternalIPExtension      `yaml:"gceExternalIp,omitempty"`
	GCEIPForwarding   GCEIPForwardingExtension    `yaml:"gceIpForwarding,omitempty"`
	GCEFirewall       GCEFirewallExtension        `yaml:"gceFirewall,omitempty"`
	GCEFirewallRange  GCEFirewallRangeExtension   `yaml:"gceFirewallRange,omitempty"`
	GCEZone           GCELocationExtension        `yaml:"gceZone,omitempty"`
	GCERegion         GCELocationExtension        `yaml:"gceRegion,omitempty"`
	IAMServiceAccount IAMServiceAccountExtension  `yaml:"iamServiceAccount,omitempty"`
}

type GCELocationExtension struct {
	AllowlistedZones   []string `yaml:"allowlistedZones,omitempty"`
	AllowlistedRegions []string `yaml:"allowlistedRegions,omitempty"`
}

type GCEMachineTypeExtension struct {
	// Minimum cpu. Used to filter the list of selectable machine types.
	MinCPU int `yaml:"minCpu,omitempty"`

	// Minimum ram. Used to filter the list of selectable machine types.
	MinRAMGB int `yaml:"minRamGb,omitempty"`

	// If true, custom machine types will not be selectable.
	// More info:
	// https://cloud.google.com/compute/docs/instances/creating-instance-with-custom-machine-type
	DisallowCustomMachineTypes bool `yaml:"disallowCustomMachineTypes,omitempty"`
}

type GCEGPUTypeExtension struct {
	MachineType string `yaml:"machineType"`
	GPUType     string `yaml:"gpuType"`
}

type GCEGPUCountExtension struct {
	// This field references another variable from the schema,
	// which must have type GCEMachineType.
	MachineTypeVariable string `yaml:"machineTypeVariable"`
}

type GCEDiskSizeExtension struct {
	// The allowable range of disk sizes depends on the disk type. This field
	// references another variable from the schema, which must have type GCEDiskType.
	DiskTypeVariable string `yaml:"diskTypeVariable"`
}

type GCENetworkExtension struct {
	// AllowSharedVpcs indicates this solution can receive
	// shared VPC selflinks (fully qualified compute links).
	AllowSharedVPCs bool `yaml:"allowSharedVpcs"`
	// Used to indicate to which machine type this network interface will be
	// attached to.
	MachineTypeVariable string `yaml:"machineTypeVariable"`
}

type ExternalIPType string

const (
	IPEphemeral ExternalIPType = "EPHEMERAL"
	IPStaic     ExternalIPType = "STATIC"
)

type GCEExternalIPExtension struct {
	// NetworkVariable is used to indicate the network variable in the schema
	// this external IP belongs to.
	NetworkVariable string `yaml:"networkVariable"`

	// Type specifies if the external IP is ephemeral or static.
	// Defaults to ephemeral if not specified.
	Type ExternalIPType `yaml:"externalIpType,omitempty"`
}

type GCEIPForwardingExtension struct {
	// NetworkVariable is used to indicate the network variable in the schema
	// this external IP belongs to.
	NetworkVariable string `yaml:"networkVariable"`
	NotConfigurable bool   `yaml:"notConfigurable"`
}

type GCEFirewallExtension struct {
	// NetworkVariable is used to indicate the network variable in the schema
	// this external IP belongs to.
	NetworkVariable string `yaml:"networkVariable"`
}

type GCEFirewallRangeExtension struct {
	// FirewallVariable is used to indicate the firewall variable with the type
	// GCEFirewall in the schema to which this firewall range belongs to.
	FirewallVariable string `yaml:"firewallVariable"`
}

type GCESubnetworkExtension struct {
	// Subnetwork variable requires a network context in order to determine the
	// set of available subnetworks. This field references another
	// variable from the schema, which must have type GCENetwork.
	NetworkVariable string `yaml:"networkVariable"`
}

type GCEGenericResourceExtension struct {
	// GCE resource type to be fetched. This field references another
	// property from the schema, which must have type GCEGenericResource.
	ResourceVariable string `yaml:"resourceVariable"`
}

type IAMServiceAccountExtension struct {
	// List of IAM roles that to  grant to a new SA, or the roles to filter
	// existing SAs with.
	Roles []string `yaml:"roles"`
}
