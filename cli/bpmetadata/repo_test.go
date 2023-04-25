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
			got := getBpRootPath(tt.path)
			if got != tt.want {
				t.Errorf("getBpRootPath() = %v, want %v", got, tt.want)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBpSubmoduleName(tt.path)
			if got != tt.want {
				t.Errorf("getBpSubmoduleName() = %v, want %v", got, tt.want)
			}
		})
	}
}
