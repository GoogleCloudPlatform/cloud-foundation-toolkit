package reports.data

import data.validator.gcp.lib as lib
import data.assets as assets

keys_report[{
    "project_id": project_id ,
    "location": location,
    "key_ring": key_ring,
    "key": key,
    "primary_version_create_time": primary_version_create_time,
    "next_rotation_time": next_rotation_time

}] {
    k := assets[_]
    count({k.asset_type} & {"cloudkms.googleapis.com/CryptoKey","google.cloud.kms.CryptoKey"}) == 1
    name_parts := split(k.resource.data.name,"/")
    project_id := name_parts[1]
    location := name_parts[3]
    key_ring := name_parts[5]
    key := name_parts[7]
    primary_version_create_time := k.resource.data.primary.createTime
    next_rotation_time := lib.get_default(k.resource.data, "nextRotationTime", "")
    
}

keys_project_level_iam_report[{
    "project_id": project_id,
    "key": key ,
    "iam_role": role,
    "iam_member": member
}] {
    k := assets[_]
    count({k.asset_type} & {"cloudkms.googleapis.com/CryptoKey","google.cloud.kms.CryptoKey"}) == 1
    key := k.resource.data.name
    name_parts := split(k.resource.data.name,"/")
    project_id := name_parts[1]

    p := assets[_]
    count({p.asset_type} & {"compute.googleapis.com/Project","google.compute.Project"}) == 1
    p.resource.data.name == project_id

    i := assets[_]
    i.name == p.resource.parent
    b := i.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    member_type := str_array[0] 
    re_match("(roles/owner|roles/editor|^roles/cloudkms)", role)
}

bucket_default_acl_report[{
    "name": b.name,
    "default_object_acl": default_object_acl_str,
    "bucket_policy_enabled": bucket_policy_enabled_str
}] {
    b := assets[_]
    count({b.asset_type} & {"storage.googleapis.com/Bucket","google.cloud.storage.Bucket"}) == 1
    default_object_acl_str := lib.is_null_str(b.resource.data.defaultObjectAcl)
    iam_configuration := lib.get_default(b, "iamConfiguration", {})
    bucket_policy_only := lib.get_default(iam_configuration, "bucketPolicyOnly", {})
    bucket_policy_enabled := lib.get_default(bucket_policy_only, "enabled", false)
    bucket_policy_enabled_str := lib.bool_to_str(bucket_policy_enabled)
}

bucket_object_lifecycle_report[{
    "name": b.name,
    "lifecycle_rule": lifecycle_rule_str
}] {
    b := assets[_]
    count({b.asset_type} & {"storage.googleapis.com/Bucket","google.cloud.storage.Bucket"}) == 1
    lifecycle_rule := lib.get_default(b.resource.data.lifecycle, "rule", [])
    lifecycle_rule_str := lib.is_null_str(lifecycle_rule)
}

bucket_location_report[{
    "name": b.name,
    "location": location
}] {
    b := assets[_]
    count({b.asset_type} & {"storage.googleapis.com/Bucket","google.cloud.storage.Bucket"}) == 1
    location := b.resource.data.location
}

dataset_no_iam_report[{
    "name": ds.name
}] {
    ds := assets[_]
    count({ds.asset_type} & {"bigquery.googleapis.com/Dataset","google.cloud.bigquery.Dataset"}) == 1

    ds_iam := assets[_]
    count({ds_iam.asset_type} & {"bigquery.googleapis.com/Dataset","google.cloud.bigquery.Dataset"}) == 1
    ds_iam.iam_policy != null
    count({ds.name} & cast_set(ds_iam[_].name)) == 0

}

dataset_project_level_iam_report[{
    "name": ds.name,
    "project": p.resource.data.name,
    "iam_role": role,
    "iam_member": member
}] {
    ds := assets[_]
    count({ds.asset_type} & {"bigquery.googleapis.com/Dataset","google.cloud.bigquery.Dataset"}) == 1

    p := assets[_]
    count({p.asset_type} & {"compute.googleapis.com/Project","google.compute.Project"}) == 1
    p.resource.parent == ds.resource.parent

    i := assets[_]
    i.name == p.resource.parent
    b := i.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    member_type := str_array[0] 
    re_match("(roles/owner|roles/editor|^roles/bigquery)", role)
}


cloud_sql_public_authorized_networks_report[{
    "name": name
}] {
    a := assets[_]
    count({a.asset_type} & {"sqladmin.googleapis.com/Instance","google.cloud.sql.Instance"}) == 1
    name := a.name
    authorized_networks := a.resource.data.settings.ipConfiguration.authorizedNetworks[_].value
    authorized_networks = "0.0.0.0/0"
}

cloud_sql_gen_report[{
    "name": name,
    "backend_type": backend_type
}] {
    a := assets[_]
    count({a.asset_type} & {"sqladmin.googleapis.com/Instance","google.cloud.sql.Instance"}) == 1
    name = a.name
    backend_type = a.resource.data.backendType
}