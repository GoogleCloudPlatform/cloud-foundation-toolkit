// Package launchpad contains libraries for reading Cloud Foundation Toolkit custom
// resource definitions and output Infrastructure as Code ready code.
//
// Supported resources can be found in generic.go
//
// == Process ==
// Launchpad divide the processing into 1) loading, 2) validating & assembling,
// 3) generating output.
//
// Input YAMLs will be organized and assembled into an GCP organization in memory.
// And generate output functions will take any GCP Org structure to output code.
//
// == References ==
// Within and across resource definitions, reference can be made to signify a target.
// Target resource can be referenced in two distinct ways, implicitly or explicitly.
//
// An implicit reference occurs if a resource is nested under another resource.
//
//   kind: Folder
//   spec:
//     id: X
//     folders:
//       - id: Y
//
// Folder Y have an implicit reference of Folder X as a parent.
//
// An explicit reference occurs if a resource specified referenced type and id.
//
//   kind: Folder
//   spec:
//     id: Y
//     parentRef:
//       type: Folder
//       id: X
//
// Folder Y have an explicit reference of Folder X as a parent.
//
// As references can have multiple use cases, all YAML definition will use `Ref`
// as suffix for referenced fields. For example, `parentRef` to specify
// parent-child relationship such as organization/folder, folder/folder.
package launchpad