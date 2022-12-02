package util

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

const (
	testContentPath = "../testdata/bpmetadata"
)

func TestTFDirectories(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr bool
	}{
		{
			name: "multiple directores",
			path: "content/examples",
			want: []string{
				"../testdata/bpmetadata/content/examples/acm/acm-terraform-blog-part1/terraform",
				"../testdata/bpmetadata/content/examples/acm/acm-terraform-blog-part2",
				"../testdata/bpmetadata/content/examples/simple_regional",
				"../testdata/bpmetadata/content/examples/simple_regional_beta",
			},
		},
		{
			name: "single directory",
			path: "content/examples/simple_regional_beta",
			want: []string{
				"../testdata/bpmetadata/content/examples/simple_regional_beta",
			},
		},
		{
			name:    "single directory",
			path:    "content/no_directory",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := path.Join(testContentPath, tt.path)
			got, err := WalkTerraformDirs(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalkTerraformDirs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got, tt.want)
		})
	}
}
