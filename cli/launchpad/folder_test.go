package launchpad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestFolder generates a fully formed folder for testing usage.
func newTestFolder(id string, name string, parType crdKind, parId string, subFolderIds []string) folderYAML {
	var subdir []*folderSpecYAML
	for _, sId := range subFolderIds {
		subdir = append(subdir, &folderSpecYAML{Id: sId, DisplayName: sId + "folder"})
	}
	return folderYAML{Spec: folderSpecYAML{
		Id: id, DisplayName: name,
		ParentRef:      referenceYAML{parType.String(), parId},
		SubFolderSpecs: subdir,
	}}
}

func TestFolders_add(t *testing.T) {
	f1, f2 := &folderYAML{Spec: folderSpecYAML{Id: "1"}}, &folderYAML{Spec: folderSpecYAML{Id: "2"}}
	folders1 := folders([]*folderYAML{f1})
	folders12 := folders([]*folderYAML{f1, f2})

	var testCases = []struct {
		name   string
		input  []*folderYAML
		output folders
	}{
		{"add_once", []*folderYAML{f1}, folders1},
		{"add_twice", []*folderYAML{f1, f2}, folders12},
		{"add_same_twice", []*folderYAML{f1, f2, f1}, folders12},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := folders{}
			for _, newF := range tc.input {
				_ = fs.add(newF) // silently ignore existing resource
			}
			assert.Equal(t, tc.output, fs, "expected folders to be the same")
		})
	}
}

func TestFolderYAML_validate(t *testing.T) {
	orgRef := referenceYAML{Organization.String(), "12345"}
	var testCases = []struct {
		name   string
		input  *folderYAML
		output error
	}{
		{"missing_id", &folderYAML{Spec: folderSpecYAML{Id: ""}}, errMissingRequiredField},
		{"no_parents", &folderYAML{Spec: folderSpecYAML{Id: "f1"}}, errInvalidParent},
		{
			"invalid_parent",
			&folderYAML{Spec: folderSpecYAML{Id: "f1", ParentRef: referenceYAML{CloudFoundation.String(), "dummy"}}},
			errInvalidParent,
		},
		{"folderName_too_short", &folderYAML{Spec: folderSpecYAML{Id: "f1", DisplayName: "12", ParentRef: orgRef}}, errValidationFailed},
		{"folderName_too_long", &folderYAML{Spec: folderSpecYAML{Id: "f1", DisplayName: "1234567890123456789012345678901", ParentRef: orgRef}}, errValidationFailed},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.validate()
			assert.Equal(t, tc.output, err)
		})
	}
}

func TestFolderYAML_validateSubFolder(t *testing.T) {
	f := newTestFolder("f1", "fold1", Organization, "12345", []string{"sf1", "sf2"})
	assert.Len(t, f.subFolders, 0, "no default subfolder should exist")
	err := f.validate()
	assert.Nil(t, err, "does not expect validation failure")

	sfs := f.subFolders // validated wrapped folders
	assert.Len(t, sfs, 2, "expected to have parsed folders")
	var buff []string
	for _, sf := range sfs {
		assert.Equal(t, Folder, sf.Spec.ParentRef.TargetType())
		assert.Equal(t, "f1", sf.Spec.ParentRef.TargetId)
		buff = append(buff, sf.Spec.Id)
	}
	assert.Equal(t, []string{"sf1", "sf2"}, buff, "expected to have samp ids")
}

func TestFolderYAML_addToOrg(t *testing.T) {
	// ORG <- f1 <- (sf1, sf2)
	f := newTestFolder("f1", "fold1", Organization, "12345", []string{"sf1", "sf2"})

	err := f.validate() // validation triggers subFolder population
	assert.Nil(t, err, "folder validation should pass")

	ao := newAssembledOrg()
	err = f.addToOrg(ao)
	assert.Nil(t, err, "folder should be now registered")

	// verify f1 exist and valid pointer
	f1, found := ao.resourceMap[f.resId()]
	assert.True(t, found, "f1 is not registered")
	assert.Equal(t, &f, f1.yaml, "f1 registered yaml is not the same")
	assert.Len(t, f1.inRefs, 2, "both sf1, sf2 should reference f1")

	org, found := ao.resourceMap[Organization.String()+".12345"]
	assert.True(t, found, "org should be populated")
	assert.Len(t, org.inRefs, 1, "f1 should point to org")

	// verify sub-folder sf1 and sf2 are also registered
	_, found = ao.resourceMap[Folder.String()+".sf1"]
	assert.True(t, found, "sf1 is not registered")
	_, found = ao.resourceMap[Folder.String()+".sf2"]
	assert.True(t, found, "sf2 is not registered")
}

// TestFolderYAML_resolveReferences ensure original YAML is not modified, but subFolder is updated.
func TestFolderYAML_resolveReferences(t *testing.T) {
	f1 := newTestFolder("f1", "fold1", Organization, "12345", []string{"sf1", "sf2"})

	err := f1.validate() // Trigger population of subFolder
	assert.Nil(t, err, "folder validation should pass")
	assert.Len(t, f1.Spec.SubFolderSpecs, 2, "subFolder YAML should be be 2 based on YAML")
	assert.Len(t, f1.subFolders, 2, "subFolder should be populate to 2")

	// new sf3 should be a f1 subFolder
	sf3 := newTestFolder("sf3", "sf3 folder", Folder, "f1", []string{})
	err = f1.resolveReferences([]resourceHandler{&sf3})
	assert.Nil(t, err, "resolve reference should pass")
	assert.Len(t, f1.Spec.SubFolderSpecs, 2, "subFolder YAML should NOT be modified")
	assert.Len(t, f1.subFolders, 3, "subFolder should be increase to 3")

	// double adding should not have any effect
	err = f1.resolveReferences([]resourceHandler{&sf3})
	assert.Nil(t, err, "resolve reference should pass")
	assert.Len(t, f1.Spec.SubFolderSpecs, 2, "subFolder YAML should NOT be modified")
	assert.Len(t, f1.subFolders, 3, "subFolder should be remain as 3")

	// new sf4 should NOT be a f1 subfolder
	sf4 := newTestFolder("sf4", "sf4 folder", Folder, "f100", []string{})
	err = f1.resolveReferences([]resourceHandler{&sf4})
	assert.Equal(t, errInvalidParent, err, "mismatch id")

	// new org1 should NOT be a f1 subfolder
	org1 := orgYAML{Spec: orgSpecYAML{Id: "dummy"}}
	err = f1.resolveReferences([]resourceHandler{&org1})
	assert.Equal(t, errInvalidInput, err, "impossible reference")
}

func TestFolderYAMLIntegration(t *testing.T) {
	var testCases = []struct {
		name              string
		expectedRelations []testFolderRelation
		inputYAMLs        []string
		assembleFail      bool
	}{{
		"one_folder_under_org",
		[]testFolderRelation{{"group1", true, testOrgId}},
		[]string{"../testdata/launchpad/folder/folder_1.yaml"},
		false,
	}, {
		"nested_folder_out_of_order",
		[]testFolderRelation{
			{"group1", true, testOrgId},
			{"group1_2", false, "group1"},
			{"group1_2_1", false, "group1_2"}, // Ordering on purpose
			{"group1_2_3", false, "group1_2"},
		},
		[]string{"../testdata/launchpad/folder/folder_1.yaml", "../testdata/launchpad/folder/folder_12.yaml"},
		false,
	}, {
		"deep_nested",
		[]testFolderRelation{
			{"group1", true, testOrgId},
			{"group11", false, "group1"},
			{"group12", false, "group1"},
			{"group123", false, "group12"},
		},
		[]string{"../testdata/launchpad/folder/folder_1_nested.yaml"},
		false,
	}, {
		"conflict_folder_def",
		[]testFolderRelation{},
		[]string{"../testdata/launchpad/folder/folder_1.yaml", "../testdata/launchpad/folder/folder_1_nested.yaml"},
		true,
	}, {
		"conflict_nested_folder_def",
		[]testFolderRelation{},
		[]string{"../testdata/launchpad/folder/org_2_nested_conflict.yaml"},
		true,
	}, {
		"parent_not_found",
		[]testFolderRelation{},
		[]string{"../testdata/launchpad/folder/folder_12.yaml"},
		true,
	}, {
		"org_crd_subfolder",
		[]testFolderRelation{{"group1", true, testOrgId}},
		[]string{"../testdata/launchpad/folder/org_1.yaml"},
		false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := newTestOrg("")
			addTestRelations(ao, tc.expectedRelations)

			resources := loadResources(tc.inputYAMLs)
			if tc.assembleFail {
				assert.Panics(t, func() {
					assembleResourcesToOrg(resources)
				}, "panic was expected for undefined reference")
			} else {
				assembled := assembleResourcesToOrg(resources)
				assert.Equal(t, ao.String(), assembled.String())
			}
		})
	}
}
