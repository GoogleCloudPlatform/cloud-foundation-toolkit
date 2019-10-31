package launchpad

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testdataPath string = "../testdata/launchpad"

func testdataFiles(relativeFilepath ...string) []string {
	var ret []string
	for _, f := range relativeFilepath {
		ret = append(ret, fmt.Sprintf("%s/%s", testdataPath, f))
	}
	return ret
}

func testFolder(id string, name string, parentType crdKind, parentId string) *folderYAML {
	return &folderYAML{
		commonConfigYAML{APIVersion: "cft.dev/v1alpha1", Kind: "Folder"},
		folderSpecYAML{
			Id:          id,
			DisplayName: name,
			ParentRef:   parentRefYAML{parentType, parentId},
		},
	}
}

func TestFolderLoad(tt *testing.T) {
	fG1 := testFolder("group1", "Group 1", "Organization", "12345678")
	fG12 := testFolder("group1_2", "Group 1-2", "Folder", "group1")
	fG121 := testFolder("group1_2_1", "Group 1-2-1", "Folder", "group1_2")

	// Folder Group 1-2 under Org directly
	fG12Org := *fG12
	fG12Org.Spec.ParentRef.ParentId = "12345678"
	fG12Org.Spec.ParentRef.ParentType = "Organization"

	// Expected relationship
	nestG12 := *fG12
	nestG12.Spec.Folders = append(nestG12.Spec.Folders, &fG121.Spec)
	nestG1 := *fG1
	nestG1.Spec.Folders = append(nestG1.Spec.Folders, &nestG12.Spec)

	var testCases = []struct {
		name         string
		loadErr      error
		filePatterns []string
		YAMLs        map[string]*folderYAML
	}{{
		"no_valid_files",
		errors.New("no valid YAML files given"),
		[]string{"dummy.txt"},
		make(map[string]*folderYAML),
	}, {
		"multiple_files_construct_to_full_org",
		errors.New("test"),
		testdataFiles("folder/folder_1.yaml", "folder/folder_12.yaml", "folder/folder_121.yaml"),
		map[string]*folderYAML{"group1_2": &nestG12, "group1_2_1": fG121, "group1": &nestG1},
	}, {
		"same_document_multiple_load",
		nil,
		testdataFiles("folder/folder_1.yaml", "folder/folder_1.yaml"),
		map[string]*folderYAML{"group1": fG1},
		//}, {
		// TODO (wengm@) support multi document in a file
		//	"one_file_multiple_objects",
		//	nil,
		//	testdataFiles("folder/folder_21_212.yaml"),
		//	map[string]*folderYAML{
		//		"group2_1": folderG21,
		//		"group2_1_2": folderG212,
		//	},
	}, {
		"nested_folders",
		nil,
		testdataFiles("folder/folder_1.yaml", "folder/folder_12_121.yaml"),
		map[string]*folderYAML{"group1_2": &nestG12, "group1_2_1": fG121, "group1": &nestG1},
	}, {
		"folder_under_cloudfoundation_CRD",
		nil,
		testdataFiles("folder/cft_1.yaml"),
		map[string]*folderYAML{"group1_2": &fG12Org},
	}, {
		"folder_under_org_CRD",
		nil,
		testdataFiles("folder/org_1.yaml"),
		map[string]*folderYAML{"group1_2": &fG12Org},
	}}

	for _, tc := range testCases {
		tt.Run(tc.name, func(t *testing.T) {
			evalYAMLs := &gState.evaluated.folders.YAMLs // Evaluated YAMLs

			defer func() { gState.init() }()
			err := loadAllYAMLs(tc.filePatterns)

			if err != nil && tc.loadErr != nil && err.Error() != tc.loadErr.Error() {
				t.Fatalf("mismatch load error, want: %v, got: %v", tc.loadErr, err)
			}
			if err != nil {
				return
			}
			// gState loaded
			assert.Equal(t, len(tc.YAMLs), len(*evalYAMLs), "mismatched folder count")
			for k, v := range tc.YAMLs {
				if evaled, ok := (*evalYAMLs)[k]; ok {
					assert.Equal(t, v.Spec.Id, evaled.Id)
					assert.Equal(t, v.Spec.Folders, evaled.Folders)
					assert.Equal(t, v.Spec.ParentRef, evaled.ParentRef)
					assert.Equal(t, v.Spec.DisplayName, evaled.DisplayName)
				}
			}
		})
	}
}
