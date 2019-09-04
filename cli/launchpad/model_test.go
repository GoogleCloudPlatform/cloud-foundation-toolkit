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

func TestFolderLoad(tt *testing.T) {
	folderG12 := &folderYAML{
		commonConfigYAML{
			APIVersion: "cft.dev/v1alpha1",
			Kind:       "Folder",
		},
		folderSpecYAML{
			Id:          "group1_2",
			DisplayName: "Group 1-2",
			ParentRef: parentRefYAML{
				"Folder",
				"group1",
			},
		},
	}
	//folderG21 := &folderYAML{
	//	commonConfigYAML{
	//		APIVersion: "cft.dev/v1alpha1",
	//		Kind:       "Folder",
	//	},
	//	folderSpecYAML{
	//		Id:          "group2_1",
	//		DisplayName: "Group 2-1",
	//		ParentRef: parentRefYAML{
	//			"folder",
	//			"group2",
	//		},
	//	},
	//}
	folderG121Spec := folderSpecYAML{
		Id:          "group1_2_1",
		DisplayName: "Group 1-2-1",
		ParentRef: parentRefYAML{
			"Folder",
			"group1_2",
		},
	}
	folderG121 := &folderYAML{
		commonConfigYAML{
			APIVersion: "cft.dev/v1alpha1",
			Kind:       "Folder",
		},
		folderG121Spec,
	}
	//folderG212 := &folderYAML{
	//	commonConfigYAML{
	//		APIVersion: "cft.dev/v1alpha1",
	//		Kind:       "Folder",
	//	},
	//	folderSpecYAML{
	//		Id:          "group2_1_2",
	//		DisplayName: "Group 2-1-2",
	//		ParentRef: parentRefYAML{
	//			"folder",
	//			"group2_1",
	//		},
	//	},
	//}
	folderG12G121 := *folderG12
	folderG12G121.Spec.Folders = append(folderG12G121.Spec.Folders, &folderG121Spec)
	folderOrgG12 := *folderG12
	folderOrgG12.Spec.ParentRef.ParentType = KindOrganization
	folderOrgG12.Spec.ParentRef.ParentId = "12345678"

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
		"multiple_files_of_one_document",
		errors.New("test"),
		testdataFiles("folder/folder_12.yaml", "folder/folder_121.yaml"),
		map[string]*folderYAML{"group1_2_1": folderG121, "group1_2": folderG12},
	}, {
		"same_document_multiple_load",
		nil,
		testdataFiles("folder/folder_12.yaml", "folder/folder_12.yaml"),
		map[string]*folderYAML{"group1_2": folderG12},
		// TODO (wengm@) support multi document in a file
		//}, {
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
		testdataFiles("folder/folder_12_121.yaml"),
		map[string]*folderYAML{"group1_2": &folderG12G121, "group1_2_1": folderG121},
	}, {
		"folder_under_cloudfoundation_CRD",
		nil,
		testdataFiles("folder/cft_1.yaml"),
		map[string]*folderYAML{"group1_2": &folderOrgG12},
	}, {
		"folder_under_org_CRD",
		nil,
		testdataFiles("folder/org_1.yaml"),
		map[string]*folderYAML{"group1_2": &folderOrgG12},
	}}

	for _, tc := range testCases {
		tt.Run(tc.name, func(t *testing.T) {
			evalYAMLs := &gState.evaluated.folders.YAMLs // Evaluated YAMLs

			defer func() { *evalYAMLs = make(map[string]*folderSpecYAML) }()
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
