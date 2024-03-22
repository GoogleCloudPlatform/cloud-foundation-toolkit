package util

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
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
			name: "multiple directories",
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

func TestFindFilesWithPattern(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		pattern   string
		skipPaths []string
		want      []string
		wantErr   bool
	}{
		{
			name:    "pattern for metadata files",
			path:    "",
			pattern: `^metadata(?:.display)?.yaml$`,
			want: []string{
				"../testdata/bpmetadata/content/examples/acm/acm-terraform-blog-part1/terraform/metadata.yaml",
				"../testdata/bpmetadata/content/examples/acm/metadata.display.yaml",
				"../testdata/bpmetadata/content/examples/acm/metadata.yaml",
				"../testdata/bpmetadata/content/examples/simple_regional/metadata.yaml",
			},
		},
		{
			name:    "pattern for tf files",
			path:    "content/examples/simple_regional",
			pattern: `.+.tf$`,
			want: []string{
				"../testdata/bpmetadata/content/examples/simple_regional/main.tf",
				"../testdata/bpmetadata/content/examples/simple_regional/modules/submodule-01/main.tf",
			},
		},
		{
			name: "pattern for tf files skipping a path",
			path: "content/examples",
			skipPaths: []string{
				"examples/acm",
			},
			pattern: `.+.tf$`,
			want: []string{
				"../testdata/bpmetadata/content/examples/simple_regional/main.tf",
				"../testdata/bpmetadata/content/examples/simple_regional/modules/submodule-01/main.tf",
				"../testdata/bpmetadata/content/examples/simple_regional_beta/main.tf",
				"../testdata/bpmetadata/content/examples/simple_regional_beta/variables.tf",
			},
		},
		{
			name: "pattern for tf files skipping multiple paths",
			path: "content/examples",
			skipPaths: []string{
				"examples/acm",
				"examples/simple_regional_beta",
			},
			pattern: `.+.tf$`,
			want: []string{
				"../testdata/bpmetadata/content/examples/simple_regional/main.tf",
				"../testdata/bpmetadata/content/examples/simple_regional/modules/submodule-01/main.tf",
			},
		},
		{
			name:    "pattern for avoiding non-metadata yaml files",
			path:    "schema",
			pattern: `^metadata(?:.display)?.yaml$`,
			want:    []string{},
		},
		{
			name:    "invalid pattern",
			pattern: `*.txt`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := path.Join(testContentPath, tt.path)
			got, err := FindFilesWithPattern(path, tt.pattern, tt.skipPaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFilesWithPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
