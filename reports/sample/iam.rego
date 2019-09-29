package reports.iam

import data.validator.gcp.lib as lib
import data.assets as assets

service_accounts_report[{
    "name": name, 
    "email": email
}] {
    a := assets[_]
    count({a.asset_type} & {"iam.googleapis.com/ServiceAccount","google.iam.ServiceAccount"}) == 1
    name := a.name
    email := a.resource.data.email
}

bindings_report[{
    "name": name,
    "role": role,
    "member": member,
    "member_type": member_type
}]{
    a := assets[_]
    name := a.name
    b := a.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    member_type := str_array[0] 
}

bindings_sa_report[{
    "name": name,
    "role": role,
    "member": member
}]{
    a := assets[_]
    name := a.name
    b := a.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    str_array[0] = "serviceAccount"
}

bindings_group_report[{
    "name": name,
    "role": role,
    "member": member
}]{
    a := assets[_]
    name := a.name
    b := a.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    str_array[0] = "group"
}

bindings_user_report[{
    "name": name,
    "role": role,
    "member": member
}]{
    a := assets[_]
    name := a.name
    b := a.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    str_array[0] = "user"
}

bindings_special_report[{
    "name": name,
    "role": role,
    "member": member
}]{
    a := assets[_]
    name := a.name
    b := a.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    re_match("(allUsers|allAuthenticatedUsers|domain)", str_array[0])
}

bindings_primitive_report[{
    "name": name,
    "role": role,
    "member": member,
    "member_type": member_type
}]{
    a := assets[_]
    name := a.name
    b := a.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    member_type := str_array[0] 
    re_match("(roles/owner|roles/editor|roles/viewer)", role)
}

bindings_networkuser_report[{
    "name": name,
    "role": role,
    "member": member,
    "member_type": member_type
}]{
    a := assets[_]
    name := a.name
    b := a.iam_policy.bindings[_] 
    role := b.role
    member := b.members[_]
    str_array := split(member, ":")
    member_type := str_array[0] 
    re_match("(roles/compute.networkUser)", role)
}

audit_logs_report[{
	"name": name,
    "type": a.asset_type,
    "service": config.service,
    "log_type":log_type
}]{
	asset_types := {
		"cloudresourcemanager.googleapis.com/Organization",
		"cloudresourcemanager.googleapis.com/Folder",
		"cloudresourcemanager.googleapis.com/Project",
	}
	a := assets[_]
    name := a.name
    a.asset_type = asset_types[_] 
    configs := lib.get_default(a.iam_policy, "audit_configs", {})
	config := configs[_]
	log_type_map := {
		1: "ADMIN_READ",
		2: "DATA_WRITE",
		3: "DATA_READ"
	}
    log_type := log_type_map[config.audit_log_configs[_].log_type]

}