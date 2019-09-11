// ----------------------------------------------------------------------------
//
//     This file is copied here by Magic Modules. The code for building up a
//     sql database instance object is copied from the manually implemented
//     provider file:
//     third_party/terraform/resources/resource_sql_database_instance.go.erb.go
//
// ----------------------------------------------------------------------------
package google

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func GetSQLDatabaseInstanceCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	name, err := assetName(d, config, "//cloudsql.googleapis.com/projects/{{project}}/instances/{{name}}")
	if err != nil {
		return Asset{}, err
	}
	if obj, err := GetSQLDatabaseInstanceApiObject(d, config); err == nil {
		return Asset{
			Name: name,
			Type: "sqladmin.googleapis.com/Instance",
			Resource: &AssetResource{
				Version:              "v1beta4",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/sqladmin/v1beta4/rest",
				DiscoveryName:        "DatabaseInstance",
				Data:                 obj,
			},
		}, nil
	} else {
		return Asset{}, err
	}
}

func GetSQLDatabaseInstanceApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		name = resource.UniqueId()
	}

	instance := &sqladmin.DatabaseInstance{
		Project:              project,
		Name:                 name,
		Region:               region,
		Settings:             expandSqlDatabaseInstanceSettings(d.Get("settings").([]interface{}), !isFirstGen(d)),
		DatabaseVersion:      d.Get("database_version").(string),
		MasterInstanceName:   d.Get("master_instance_name").(string),
		ReplicaConfiguration: expandReplicaConfiguration(d.Get("replica_configuration").([]interface{})),
	}

	return jsonMap(instance)
}

// Detects whether a database is 1st Generation by inspecting the tier name
func isFirstGen(d TerraformResourceData) bool {
	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})
	tier := settings["tier"].(string)

	// 1st Generation databases have tiers like 'D0', as opposed to 2nd Generation which are
	// prefixed with 'db'
	return !regexp.MustCompile("db*").Match([]byte(tier))
}

func expandSqlDatabaseInstanceSettings(configured []interface{}, secondGen bool) *sqladmin.Settings {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_settings := configured[0].(map[string]interface{})
	settings := &sqladmin.Settings{
		// Version is unset in Create but is set during update
		SettingsVersion:             int64(_settings["version"].(int)),
		Tier:                        _settings["tier"].(string),
		ForceSendFields:             []string{"StorageAutoResize"},
		ActivationPolicy:            _settings["activation_policy"].(string),
		AvailabilityType:            _settings["availability_type"].(string),
		CrashSafeReplicationEnabled: _settings["crash_safe_replication"].(bool),
		DataDiskSizeGb:              int64(_settings["disk_size"].(int)),
		DataDiskType:                _settings["disk_type"].(string),
		PricingPlan:                 _settings["pricing_plan"].(string),
		ReplicationType:             _settings["replication_type"].(string),
		UserLabels:                  convertStringMap(_settings["user_labels"].(map[string]interface{})),
		BackupConfiguration:         expandBackupConfiguration(_settings["backup_configuration"].([]interface{})),
		DatabaseFlags:               expandDatabaseFlags(_settings["database_flags"].([]interface{})),
		AuthorizedGaeApplications:   expandAuthorizedGaeApplications(_settings["authorized_gae_applications"].([]interface{})),
		IpConfiguration:             expandIpConfiguration(_settings["ip_configuration"].([]interface{})),
		LocationPreference:          expandLocationPreference(_settings["location_preference"].([]interface{})),
		MaintenanceWindow:           expandMaintenanceWindow(_settings["maintenance_window"].([]interface{})),
	}

	// 1st Generation instances don't support the disk_autoresize parameter
	// and it defaults to true - so we shouldn't set it if this is first gen
	if secondGen {
		settings.StorageAutoResize = googleapi.Bool(_settings["disk_autoresize"].(bool))
	}

	return settings
}

func expandReplicaConfiguration(configured []interface{}) *sqladmin.ReplicaConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_replicaConfiguration := configured[0].(map[string]interface{})
	return &sqladmin.ReplicaConfiguration{
		FailoverTarget: _replicaConfiguration["failover_target"].(bool),

		// MysqlReplicaConfiguration has been flattened in the TF schema, so
		// we'll keep it flat here instead of another expand method.
		MysqlReplicaConfiguration: &sqladmin.MySqlReplicaConfiguration{
			CaCertificate:           _replicaConfiguration["ca_certificate"].(string),
			ClientCertificate:       _replicaConfiguration["client_certificate"].(string),
			ClientKey:               _replicaConfiguration["client_key"].(string),
			ConnectRetryInterval:    int64(_replicaConfiguration["connect_retry_interval"].(int)),
			DumpFilePath:            _replicaConfiguration["dump_file_path"].(string),
			MasterHeartbeatPeriod:   int64(_replicaConfiguration["master_heartbeat_period"].(int)),
			Password:                _replicaConfiguration["password"].(string),
			SslCipher:               _replicaConfiguration["ssl_cipher"].(string),
			Username:                _replicaConfiguration["username"].(string),
			VerifyServerCertificate: _replicaConfiguration["verify_server_certificate"].(bool),
		},
	}
}

func expandMaintenanceWindow(configured []interface{}) *sqladmin.MaintenanceWindow {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	window := configured[0].(map[string]interface{})
	return &sqladmin.MaintenanceWindow{
		Day:             int64(window["day"].(int)),
		Hour:            int64(window["hour"].(int)),
		UpdateTrack:     window["update_track"].(string),
		ForceSendFields: []string{"Hour"},
	}
}

func expandLocationPreference(configured []interface{}) *sqladmin.LocationPreference {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_locationPreference := configured[0].(map[string]interface{})
	return &sqladmin.LocationPreference{
		FollowGaeApplication: _locationPreference["follow_gae_application"].(string),
		Zone:                 _locationPreference["zone"].(string),
	}
}

func expandIpConfiguration(configured []interface{}) *sqladmin.IpConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_ipConfiguration := configured[0].(map[string]interface{})

	return &sqladmin.IpConfiguration{
		Ipv4Enabled:        _ipConfiguration["ipv4_enabled"].(bool),
		RequireSsl:         _ipConfiguration["require_ssl"].(bool),
		PrivateNetwork:     _ipConfiguration["private_network"].(string),
		AuthorizedNetworks: expandAuthorizedNetworks(_ipConfiguration["authorized_networks"].(*schema.Set).List()),
		ForceSendFields:    []string{"Ipv4Enabled", "RequireSsl"},
	}
}
func expandAuthorizedNetworks(configured []interface{}) []*sqladmin.AclEntry {
	an := make([]*sqladmin.AclEntry, 0, len(configured))
	for _, _acl := range configured {
		_entry := _acl.(map[string]interface{})
		an = append(an, &sqladmin.AclEntry{
			ExpirationTime: _entry["expiration_time"].(string),
			Name:           _entry["name"].(string),
			Value:          _entry["value"].(string),
		})
	}

	return an
}

func expandAuthorizedGaeApplications(configured []interface{}) []string {
	aga := make([]string, 0, len(configured))
	for _, app := range configured {
		aga = append(aga, app.(string))
	}
	return aga
}

func expandDatabaseFlags(configured []interface{}) []*sqladmin.DatabaseFlags {
	databaseFlags := make([]*sqladmin.DatabaseFlags, 0, len(configured))
	for _, _flag := range configured {
		_entry := _flag.(map[string]interface{})

		databaseFlags = append(databaseFlags, &sqladmin.DatabaseFlags{
			Name:  _entry["name"].(string),
			Value: _entry["value"].(string),
		})
	}
	return databaseFlags
}

func expandBackupConfiguration(configured []interface{}) *sqladmin.BackupConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_backupConfiguration := configured[0].(map[string]interface{})
	return &sqladmin.BackupConfiguration{
		BinaryLogEnabled: _backupConfiguration["binary_log_enabled"].(bool),
		Enabled:          _backupConfiguration["enabled"].(bool),
		StartTime:        _backupConfiguration["start_time"].(string),
	}
}
