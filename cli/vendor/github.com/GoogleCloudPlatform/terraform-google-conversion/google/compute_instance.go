// ----------------------------------------------------------------------------
//
//     This file is copied here by Magic Modules. The code for building up a
//     compute instance object is copied from the manually implemented
//     provider file:
//     third_party/terraform/resources/resource_compute_instance.go
//
// ----------------------------------------------------------------------------
package google

import (
	"fmt"
	"log"
	"strings"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func GetComputeInstanceCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	name, err := assetName(d, config, "//compute.googleapis.com/projects/{{project}}/zones/{{zone}}/instances/{{name}}")
	if err != nil {
		return Asset{}, err
	}
	if obj, err := GetComputeInstanceApiObject(d, config); err == nil {
		return Asset{
			Name: name,
			Type: "compute.googleapis.com/Instance",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Instance",
				Data:                 obj,
			},
		}, nil
	} else {
		return Asset{}, err
	}
}

func GetComputeInstanceApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	// Get the zone
	z, err := getZone(d, config)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Loading zone: %s", z)
	zone, err := config.clientCompute.Zones.Get(
		project, z).Do()
	if err != nil {
		return nil, fmt.Errorf("Error loading zone '%s': %s", z, err)
	}

	instance, err := expandComputeInstance(project, zone, d, config)
	if err != nil {
		return nil, err
	}

	return jsonMap(instance)
}

func expandComputeInstance(project string, zone *compute.Zone, d TerraformResourceData, config *Config) (*computeBeta.Instance, error) {
	// Get the machine type
	var machineTypeUrl string
	if mt, ok := d.GetOk("machine_type"); ok {
		log.Printf("[DEBUG] Loading machine type: %s", mt.(string))
		machineType, err := config.clientCompute.MachineTypes.Get(
			project, zone.Name, mt.(string)).Do()
		if err != nil {
			return nil, fmt.Errorf(
				"Error loading machine type: %s",
				err)
		}
		machineTypeUrl = machineType.SelfLink
	}

	// Build up the list of disks

	disks := []*computeBeta.AttachedDisk{}
	if _, hasBootDisk := d.GetOk("boot_disk"); hasBootDisk {
		bootDisk, err := expandBootDisk(d, config, zone, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, bootDisk)
	}

	if _, hasScratchDisk := d.GetOk("scratch_disk"); hasScratchDisk {
		scratchDisks, err := expandScratchDisks(d, config, zone, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, scratchDisks...)
	}

	attachedDisksCount := d.Get("attached_disk.#").(int)

	for i := 0; i < attachedDisksCount; i++ {
		diskConfig := d.Get(fmt.Sprintf("attached_disk.%d", i)).(map[string]interface{})
		disk, err := expandAttachedDisk(diskConfig, d, config)
		if err != nil {
			return nil, err
		}

		disks = append(disks, disk)
	}

	sch := d.Get("scheduling").([]interface{})
	var scheduling *computeBeta.Scheduling
	if len(sch) == 0 {
		// TF doesn't do anything about defaults inside of nested objects, so if
		// scheduling hasn't been set, then send it with its default values.
		scheduling = &computeBeta.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}
	} else {
		prefix := "scheduling.0"
		scheduling = &computeBeta.Scheduling{
			AutomaticRestart:  googleapi.Bool(d.Get(prefix + ".automatic_restart").(bool)),
			Preemptible:       d.Get(prefix + ".preemptible").(bool),
			OnHostMaintenance: d.Get(prefix + ".on_host_maintenance").(string),
			ForceSendFields:   []string{"AutomaticRestart", "Preemptible"},
		}
	}

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating metadata: %s", err)
	}

	networkInterfaces, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network interfaces: %s", err)
	}

	accels, err := expandInstanceGuestAccelerators(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating guest accelerators: %s", err)
	}

	// Create the instance information
	return &computeBeta.Instance{
		CanIpForward:       d.Get("can_ip_forward").(bool),
		Description:        d.Get("description").(string),
		Disks:              disks,
		MachineType:        machineTypeUrl,
		Metadata:           metadata,
		Name:               d.Get("name").(string),
		NetworkInterfaces:  networkInterfaces,
		Tags:               resourceInstanceTags(d),
		Labels:             expandLabels(d),
		ServiceAccounts:    expandServiceAccounts(d.Get("service_account").([]interface{})),
		GuestAccelerators:  accels,
		MinCpuPlatform:     d.Get("min_cpu_platform").(string),
		Scheduling:         scheduling,
		DeletionProtection: d.Get("deletion_protection").(bool),
		Hostname:           d.Get("hostname").(string),
		ForceSendFields:    []string{"CanIpForward", "DeletionProtection"},
	}, nil
}

func expandAttachedDisk(diskConfig map[string]interface{}, d TerraformResourceData, meta interface{}) (*computeBeta.AttachedDisk, error) {
	config := meta.(*Config)

	s := diskConfig["source"].(string)
	var sourceLink string
	if strings.Contains(s, "regions/") {
		source, err := ParseRegionDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	} else {
		source, err := ParseDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	}

	disk := &computeBeta.AttachedDisk{
		Source: sourceLink,
	}

	if v, ok := diskConfig["mode"]; ok {
		disk.Mode = v.(string)
	}

	if v, ok := diskConfig["device_name"]; ok {
		disk.DeviceName = v.(string)
	}

	if v, ok := diskConfig["disk_encryption_key_raw"]; ok {
		disk.DiskEncryptionKey = &computeBeta.CustomerEncryptionKey{
			RawKey: v.(string),
		}
	}
	return disk, nil
}

// See comment on expandInstanceTemplateGuestAccelerators regarding why this
// code is duplicated.
func expandInstanceGuestAccelerators(d TerraformResourceData, config *Config) ([]*computeBeta.AcceleratorConfig, error) {
	configs, ok := d.GetOk("guest_accelerator")
	if !ok {
		return nil, nil
	}
	accels := configs.([]interface{})
	guestAccelerators := make([]*computeBeta.AcceleratorConfig, 0, len(accels))
	for _, raw := range accels {
		data := raw.(map[string]interface{})
		if data["count"].(int) == 0 {
			continue
		}
		at, err := ParseAcceleratorFieldValue(data["type"].(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot parse accelerator type: %v", err)
		}
		guestAccelerators = append(guestAccelerators, &computeBeta.AcceleratorConfig{
			AcceleratorCount: int64(data["count"].(int)),
			AcceleratorType:  at.RelativeLink(),
		})
	}

	return guestAccelerators, nil
}

func expandBootDisk(d TerraformResourceData, config *Config, zone *compute.Zone, project string) (*computeBeta.AttachedDisk, error) {
	disk := &computeBeta.AttachedDisk{
		AutoDelete: d.Get("boot_disk.0.auto_delete").(bool),
		Boot:       true,
	}

	if v, ok := d.GetOk("boot_disk.0.device_name"); ok {
		disk.DeviceName = v.(string)
	}

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_key_raw"); ok {
		disk.DiskEncryptionKey = &computeBeta.CustomerEncryptionKey{
			RawKey: v.(string),
		}
	}

	if v, ok := d.GetOk("boot_disk.0.source"); ok {
		source, err := ParseDiskFieldValue(v.(string), d, config)
		if err != nil {
			return nil, err
		}
		disk.Source = source.RelativeLink()
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params"); ok {
		disk.InitializeParams = &computeBeta.AttachedDiskInitializeParams{}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.size"); ok {
			disk.InitializeParams.DiskSizeGb = int64(v.(int))
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.type"); ok {
			diskTypeName := v.(string)
			diskType, err := readDiskType(config, zone, project, diskTypeName)
			if err != nil {
				return nil, fmt.Errorf("Error loading disk type '%s': %s", diskTypeName, err)
			}
			disk.InitializeParams.DiskType = diskType.SelfLink
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.image"); ok {
			imageName := v.(string)
			imageUrl, err := resolveImage(config, project, imageName)
			if err != nil {
				return nil, fmt.Errorf("Error resolving image name '%s': %s", imageName, err)
			}

			disk.InitializeParams.SourceImage = imageUrl
		}
	}

	return disk, nil
}

func expandScratchDisks(d TerraformResourceData, config *Config, zone *compute.Zone, project string) ([]*computeBeta.AttachedDisk, error) {
	diskType, err := readDiskType(config, zone, project, "local-ssd")
	if err != nil {
		return nil, fmt.Errorf("Error loading disk type 'local-ssd': %s", err)
	}

	n := d.Get("scratch_disk.#").(int)
	scratchDisks := make([]*computeBeta.AttachedDisk, 0, n)
	for i := 0; i < n; i++ {
		scratchDisks = append(scratchDisks, &computeBeta.AttachedDisk{
			AutoDelete: true,
			Type:       "SCRATCH",
			Interface:  d.Get(fmt.Sprintf("scratch_disk.%d.interface", i)).(string),
			InitializeParams: &computeBeta.AttachedDiskInitializeParams{
				DiskType: diskType.SelfLink,
			},
		})
	}

	return scratchDisks, nil
}
