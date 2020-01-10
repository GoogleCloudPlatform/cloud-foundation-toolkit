package launchpad

const testOrgId = "12345678"

// newTestOrg generates an assembledOrg with it's org registered.
func newTestOrg(orgId string) *assembledOrg {
	if orgId == "" {
		orgId = testOrgId
	}
	ao := newAssembledOrg()
	ao.org = orgYAML{Spec: orgSpecYAML{Id: orgId}}
	ao.resourceMap[ao.org.resId()] = &resource{yaml: &ao.org}
	return ao
}

type testFolderRelation struct {
	idName      string
	parentIsOrg bool
	parentId    string
}

// addTestRelations appends children relationships to resource under an initialized assemble org.
//
// addTestRelations will also initialize resources for resourceMap. his test method assumes input
// relations are purposely ordered so when looping it though the parents already exist in resourceMap.
func addTestRelations(ao *assembledOrg, testFolders []testFolderRelation) {
	for _, fdr := range testFolders {
		parentType := Folder.String()
		if fdr.parentIsOrg {
			parentType = Organization.String()
		}
		f := &folderYAML{Spec: folderSpecYAML{
			Id: fdr.idName, DisplayName: fdr.idName,
			ParentRef: referenceYAML{parentType, fdr.parentId},
		}}
		ao.resourceMap[f.resId()] = &resource{yaml: f}

		parentId := f.Spec.ParentRef.resId()
		switch parent := ao.resourceMap[parentId].yaml.(type) {
		case *orgYAML:
			parent.subFolders = append(parent.subFolders, f)
		case *folderYAML:
			parent.subFolders = append(parent.subFolders, f)
		}
		ao.resourceMap[parentId].inRefs = append(ao.resourceMap[f.Spec.ParentRef.resId()].inRefs, f)
	}
}
