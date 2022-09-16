package bpmetadata

import (
	"testing"
)

func TestGetBpRepoPath(t *testing.T) {
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
			name:    "invalid top level",
			path:    "testdata/bpmetadata/erraform-google-bp01",
			wantErr: true,
		},
		{
			name:    "invalid nested",
			path:    "testdata/bpmetadata/terraform-google-bp01/test/bp01-01",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBpPathForRepoName(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBpPathForRepoName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBpPathForRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
