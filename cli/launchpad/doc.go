// Package launchpad contains libraries for reading Cloud Foundation Toolkit custom
// resource definitions and output Infrastructure as Code ready scripts.
//
// # Supported resources can be found in generic.go
//
// All resources implements resourceHandler interface and are expected
// to have a full YAML representation with name {resource}YAML and additional
// {resource}SpecYAML to denote it's Spec.
//
// Resources can reside under another resource represented by {resource}SpecYAML list,
// however, other functions will expect {resource}YAML for actual processing.
// Implementer of a resource type should aim to track sub resources as an addition
// field, as oppose to manipulating the parsed YAML directly. For example,
//
//	  kind: Folder
//	  spec:
//	    id: X
//	    folders:
//	      - id: Y
//			 - id: Z
//
// Folder X have Y, Z folders as sub resources. In the evaluation hierarchy, folder X
// is a folderYAML representation with folderSpecYAML spec. Subdirectories Y, Z are
// represented by folderSpecYAML. During validation of folder X, sub-folder Y and
// Z will be wrapped into a fully qualified folderYAML and track alongside
// without changing the original YAML definition.
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
//	kind: Folder
//	spec:
//	  id: X
//	  folders:
//	    - id: Y
//
// Folder Y have an implicit reference of Folder X as a parent.
//
// An explicit reference occurs if a resource specified referenced type and id.
//
//	kind: Folder
//	spec:
//	  id: Y
//	  parentRef:
//	    type: Folder
//	    id: X
//
// Folder Y have an explicit reference of Folder X as a parent.
//
// As references can have multiple use cases, all YAML definition will use `Ref`
// as suffix for referenced fields. For example, `parentRef` to specify
// parent-child relationship such as organization/folder, folder/folder.
package launchpad
