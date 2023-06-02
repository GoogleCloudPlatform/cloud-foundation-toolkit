package bpmetadata

import (
	"testing"
)

func TestGetBpRootPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "simple",
			path:    "testdata/bpmetadata/terraform-google-bp01",
			want:    "testdata/bpmetadata/terraform-google-bp01",
			wantErr: false,
		},
		{
			name:    "one level nested",
			path:    "testdata/bpmetadata/terraform-google-bp01/modules/bp01-01",
			want:    "testdata/bpmetadata/terraform-google-bp01",
			wantErr: false,
		},
		{
			name:    "two level nested",
			path:    "testdata/bpmetadata/terraform-google-bp01/modules/bp01-01/subbp01-01",
			want:    "testdata/bpmetadata/terraform-google-bp01",
			wantErr: false,
		},
		{
			name:    "docker workspace root",
			path:    "workspace",
			want:    "workspace",
			wantErr: false,
		},
		{
			name:    "docker workspace submodule",
			path:    "workspace/modules/bp-01",
			want:    "workspace",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBlueprintRootPath(tt.path)
			if got != tt.want {
				t.Errorf("getBlueprintRootPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBpSubmoduleName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "simple - no submodules",
			path: "testdata/bpmetadata/terraform-google-bp01",
			want: "",
		},
		{
			name: "simple - valid submodule",
			path: "testdata/bpmetadata/terraform-google-bp01/modules/submodule-01",
			want: "submodule-01",
		},
		{
			name: "simple - invalid submodule",
			path: "testdata/bpmetadata/terraform-google-bp01/foo/submodule-01",
			want: "",
		},
		{
			name: "simple - submodule with underscores",
			path: "testdata/bpmetadata/terraform-google-bp01/modules/submodule_01",
			want: "submodule-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBpSubmoduleNameInKebabCase(tt.path)
			if got != tt.want {
				t.Errorf("getBpSubmoduleNameInKebabCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRepoDetailsFromRootBp(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		wantRepoDetails repoDetail
	}{
		{
			name: "root metadata does not exist",
			path: "../testdata/bpmetadata/content/examples/simple_regional_beta/modules/submodule-01",
			wantRepoDetails: repoDetail{
				Source: &repoSource{
					BlueprintRootPath: "../testdata/bpmetadata/content/examples/simple_regional_beta",
				},
			},
		},
		{
			name: "root metadata exists but does not have source info",
			path: "../testdata/bpmetadata/content/examples/acm/modules/submodule-01",
			wantRepoDetails: repoDetail{
				RepoName: "terraform-google-acm",
				Source: &repoSource{
					BlueprintRootPath: "../testdata/bpmetadata/content/examples/acm",
				},
			},
		},
		{
			name: "root metadata exists and has source info",
			path: "../testdata/bpmetadata/content/examples/simple_regional/modules/submodule-01",
			wantRepoDetails: repoDetail{
				RepoName: "simple-regional",
				Source: &repoSource{
					URL:               "https://github.com/GoogleCloudPlatform/simple-regional",
					SourceType:        "git",
					BlueprintRootPath: "../testdata/bpmetadata/content/examples/simple_regional",
					RepoRootPath:      "../testdata/bpmetadata/content/examples/simple_regional",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRepoDetailsFromRootBp(tt.path)
			if got.Source != nil && *got.Source != *tt.wantRepoDetails.Source {
				t.Errorf("getRepoDetailsFromRootBp() - Source = %v, want %v", *got.Source, *tt.wantRepoDetails.Source)
			}

			if got.RepoName != tt.wantRepoDetails.RepoName {
				t.Errorf("getRepoDetailsFromRootBp() - RepoName = %v, want %v", got.RepoName, tt.wantRepoDetails.RepoName)
			}
		})
	}
}
