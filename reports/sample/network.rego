package reports.network

import data.validator.gcp.lib as lib
import data.assets as assets

network_default_report[{
    "name": name
}] {
    a := assets[_]
    count({a.asset_type} & {"compute.googleapis.com/Network","google.compute.Network"}) == 1
    a.resource.data.name == "default"
    name = a.name
}


firewall_egress_deny_rules_report[{
    "name": f.name,
    "network": f.resource.data.network
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    f.resource.data.direction == "EGRESS"
    denied := lib.get_default(f.resource.data, "denied", [])
    denied != []
}

firewall_egress_deny_rules[{
    "firewall": f
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    f.resource.data.direction == "EGRESS"
    denied := lib.get_default(f.resource.data, "denied", [])
    denied != []
}


firewall_ingress_allow_wide_rules_report[{
    "name": f.name,
    "network": f.resource.data.network,
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    f.resource.data.direction == "INGRESS"
    allowed := lib.get_default(f.resource.data, "allowed", [])
    allowed != []
    f.resource.data.sourceRange[_] == "0.0.0.0/0"

}


firewall_ingress_deny_wide_rules_logging_report[{
    "name": f.name,
    "network": f.resource.data.network,
    "log_enabled": log_enabled
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    f.resource.data.direction == "INGRESS"
    denied := lib.get_default(f.resource.data, "denied", [])
    denied != []
    f.resource.data.sourceRange[_] == "0.0.0.0/0"
    log_enabled = lib.bool_to_str(f.resource.data.logConfig.enable)
}



firewall_ingress_allow_ssh_rules_report[{
    "name": f.name,
    "network": f.resource.data.network
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    f.resource.data.direction == "INGRESS"
    allowed := lib.get_default(f.resource.data, "allowed", [])
    allowed != []
    allowed_rule := allowed[_]
    allowed_rule.ipProtocol == "tcp"
    allowed_rule.port[_] == "22"
}

firewall_ingress_allow_rdp_rules_report[{
    "name": f.name,
    "network": f.resource.data.network
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    f.resource.data.direction == "INGRESS"
    allowed := lib.get_default(f.resource.data, "allowed", [])
    allowed != []
    allowed_rule := allowed[_]
    allowed_rule.ipProtocol == "tcp"
    allowed_rule.port[_] == "3389"
}

firewall_service_account_report[{
    "name": f.name,
    "network": f.resource.data.network,
    "target_service_account": target_service_account
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    target_service_accounts:= lib.get_default(f.resource.data, "targetServiceAccount", [])
    target_service_accounts != []
    target_service_account := target_service_accounts[_]
}

firewall_default_report[{
    "name": name
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1

    name = f.name
    target_tag := lib.get_default(f.resource.data, "targetTag", [])
    target_tag == []
    target_service_account:= lib.get_default(f.resource.data, "targetServiceAccount", [])
    target_service_account == []

    n := assets[_]
    count({n.asset_type} & {"compute.googleapis.com/Network","google.compute.Network"}) == 1
    n.resource.data.selfLink == f.resource.data.network
    n.resource.data.name == "default"

}

firewall_logging_report[{
    "name": f.name,
    "log_enabled": log_enabled
}] {
    f := assets[_]
    count({f.asset_type} & {"compute.googleapis.com/Firewall","google.compute.Firewall"}) == 1
    log_enabled = lib.bool_to_str(f.resource.data.logConfig.enable)
}

vpc_host_projects_report[{
    "project_id": project_id,
    "xpn_project_status": xpn_project_status,
    "network": network,
    "subnetwork": subnetwork
}] {
    p := assets[_]
    count({p.asset_type} & {"compute.googleapis.com/Project","google.compute.Project"}) == 1
    project_id := p.resource.data.name
    xpn_project_status := p.resource.data.xpnProjectStatus

    n := assets[_]
    count({n.asset_type} & {"compute.googleapis.com/Network","google.compute.Network"}) == 1
    n.resource.parent == p.resource.parent
    network = n.name
    subnetwork := n.resource.data.subnetwork[_]
}

vpn_tunnels_report[{
    "project_id": project_id,
    "xpn_project_status": xpn_project_status,
    "vpc_tunnel": tunnel
}] {
    p := assets[_]
    count({p.asset_type} & {"compute.googleapis.com/Project","google.compute.Project"}) == 1
    project_id := p.resource.data.name
    xpn_project_status := p.resource.data.xpnProjectStatus

    t := assets[_]
    count({t.asset_type} & {"compute.googleapis.com/VpnTunnel","google.compute.VpnTunnel"}) == 1
    t.resource.parent == p.resource.parent
    tunnel = t.resource.data.name
}

public_vms_report[{
    "project_id": project_id,
    "instance_name": instance_name,
    "external_ip": external_ip,
    "status": status
}] {
    p := assets[_]
    count({p.asset_type} & {"compute.googleapis.com/Project","google.compute.Project"}) == 1
    project_id := p.resource.data.name
    
    vm := assets[_]
    count({vm.asset_type} & {"compute.googleapis.com/Instance","google.compute.Instance"}) == 1
    vm.resource.parent == p.resource.parent
    instance_name = vm.resource.data.name
    access_config := lib.get_default(vm.resource.data.networkInterface[_], "accessConfig", [])
    external_nat_type := lib.get_default(access_config[_], "type", "")
    external_nat_type == "ONE_TO_ONE_NAT"
    external_ip := lib.get_default(access_config[_], "externalIp", "")
    status := vm.resource.data.status
}

private_vms_report[{
    "project_id": project_id,
    "instance_name": instance_name,
    "status": status
}] {
    p := assets[_]
    count({p.asset_type} & {"compute.googleapis.com/Project","google.compute.Project"}) == 1
    project_id := p.resource.data.name
    
    vm := assets[_]
    count({vm.asset_type} & {"compute.googleapis.com/Instance","google.compute.Instance"}) == 1
    vm.resource.parent == p.resource.parent
    instance_name = vm.resource.data.name
    access_config := lib.get_default(vm.resource.data.networkInterface[_], "accessConfig", [])
    external_nat_type := lib.get_default(access_config[_], "type", "")
    external_nat_type != "ONE_TO_ONE_NAT"
    status := vm.resource.data.status
}

subnet_private_google_access_report[{
    "subnetwork_name": subnetwork_name,
    "private_google_access": private_google_access_str
}] {
    sn := assets[_]
    count({sn.asset_type} & {"compute.googleapis.com/Subnetwork","google.compute.Subnetwork"}) == 1
    subnetwork_name = sn.name
    private_google_access_str := lib.bool_to_str(sn.resource.data.privateIpGoogleAccess)
}

subnet_flow_logs_report[{
    "subnetwork_name": subnetwork_name,
    "enable_flow_logs": enable_flow_logs_str
}] {
    sn := assets[_]
    count({sn.asset_type} & {"compute.googleapis.com/Subnetwork","google.compute.Subnetwork"}) == 1
    subnetwork_name = sn.name
    enable_flow_logs := lib.get_default(sn.resource.data, "enableFlowLogs", false)
    enable_flow_logs_str := lib.bool_to_str(enable_flow_logs)
}

