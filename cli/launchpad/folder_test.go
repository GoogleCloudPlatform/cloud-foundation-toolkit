package launchpad

import "testing"
import "github.com/stretchr/testify/assert"

func TestFolders_add(t *testing.T) {
	f1, f2 := &folderYAML{Spec: folderSpecYAML{Id: "1"}}, &folderYAML{Spec: folderSpecYAML{Id: "2"}}
	folders1 := folders([]*folderYAML{f1})
	folders12 := folders([]*folderYAML{f1, f2})

	var testCases = []struct {
		name   string
		input  []*folderYAML
		output folders
	}{{
		"add_once",
		[]*folderYAML{f1},
		folders1,
	}, {
		"add_twice",
		[]*folderYAML{f1, f2},
		folders12,
	}, {
		"add_same_twice",
		[]*folderYAML{f1, f2, f1},
		folders12,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := folders{}
			for _, newF := range tc.input {
				fs.add(newF)
			}
			assert.Equal(t, tc.output, fs, "expected folders to be the same")
		})
	}
	// TODO test add folder does NOT modify original YAML
}

func TestFolderYAML_validate(t *testing.T) {
	orgRef := referenceYAML{Organization.String(), "12345"}
	var testCases = []struct {
		name   string
		input  *folderYAML
		output error
	}{{
		"missing_id",
		&folderYAML{Spec: folderSpecYAML{Id: ""}},
		errMissingRequiredField,
	}, {
		"no_parents",
		&folderYAML{Spec: folderSpecYAML{Id: "f1"}},
		errInvalidParent,
	}, {
		"invalid_parent",
		&folderYAML{Spec: folderSpecYAML{Id: "f1", ParentRef: referenceYAML{CloudFoundation.String(), "dummy"}}},
		errInvalidParent,
	}, {
		"folderName_too_short",
		&folderYAML{Spec: folderSpecYAML{Id: "f1", DisplayName: "12", ParentRef: orgRef}},
		errValidationFailed,
	}, {
		"folderName_too_long",
		&folderYAML{Spec: folderSpecYAML{Id: "f1", DisplayName: "1234567890123456789012345678901", ParentRef: orgRef}},
		errValidationFailed,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.validate()
			assert.Equal(t, tc.output, err)
		})
	}
}

var testFolder = folderYAML{Spec: folderSpecYAML{
	Id:             "f1",
	DisplayName:    "fold1",
	ParentRef:      referenceYAML{Organization.String(), "12345"},
	SubFolderSpecs: []*folderSpecYAML{{Id: "sf1", DisplayName: "subf1"}, {Id: "sf2", DisplayName: "subf2"}},
}}

func TestFolderYAML_validateSubFolder(t *testing.T) {
	assert.Len(t, testFolder.subFolders, 0, "no default subfolder should exist")
	err := testFolder.validate()
	assert.Nil(t, err, "does not expect validation failure")

	sfs := testFolder.subFolders // validated wrapped folders
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
	//f := folderYAML{Spec: folderSpecYAML{
	//	Id: "f1",
	//	DisplayName: "fold1",
	//	ParentRef:
	//}}
}

func TestFolderYAML_resolveReferences(t *testing.T) {

}
