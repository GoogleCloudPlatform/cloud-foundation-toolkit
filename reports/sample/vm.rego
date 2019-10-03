package reports.vm

import data.validator.gcp.lib as lib
import data.assets as assets


disk_source_report[{
    "project_id": project_id,
    "disk": disk_name,
    "source_image": source_image,
    "source_snapshot": source_snapshot
}] {
    p := assets[_]
    count({p.asset_type} & {"compute.googleapis.com/Project","google.compute.Project"}) == 1
    project_id := p.resource.data.name
    
    d := assets[_]
    count({d.asset_type} & {"compute.googleapis.com/Disk","google.compute.Disk"}) == 1
    d.resource.parent == p.resource.parent
    disk_name := d.resource.data.name
    source_image := lib.get_default(d.resource.data, "sourceImage", "")
    source_snapshot := lib.get_default(d.resource.data, "sourceSnapshot", "")
}

service_account_report[{
    "name": name,
    "sa_email": sa_email
}] {
    vm := assets[_]
    count({vm.asset_type} & {"compute.googleapis.com/Instance","google.compute.Instance"}) == 1
    
    name := vm.name
    sa_email := vm.resource.data.serviceAccount[_].email
}
