output "network_name" { value = "${local.network_name}" }
output "network_self_link" { value = "${module.network.network_self_link}" }
output "subnet_name" { value = "${local.subnet_name}" }
output "subnet_range_pods_name" { value = "${local.subnet_range_pods_name}" }
output "subnet_range_services_name" { value = "${local.subnet_range_services_name}" }
