package bpmetadata

import (
	"path"
	"regexp"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bptestdataPath = "../testdata/bpmetadata"
)

func TestIsPathValid(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{
			name:    "valid",
			path:    "assets/icon.png",
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid",
			path:    "assets/icon2.png",
			wantErr: true,
		},
		{
			name:    "empty",
			path:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileExists(path.Join(bptestdataPath, tt.path))
			if (err != nil) != tt.wantErr {
				t.Errorf("fileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirContent(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		regex   string
		want    []*BlueprintMiscContent
		wantErr bool
	}{
		{
			name:  "valid examples",
			path:  "content/examples",
			regex: regexExamples,
			want: []*BlueprintMiscContent{
				{
					Name:     "terraform",
					Location: "examples/acm/acm-terraform-blog-part1/terraform",
				},
				{
					Name:     "acm-terraform-blog-part2",
					Location: "examples/acm/acm-terraform-blog-part2",
				},
				{
					Name:     "simple_regional",
					Location: "examples/simple_regional",
				},
				{
					Name:     "simple_regional_beta",
					Location: "examples/simple_regional_beta",
				},
			},
			wantErr: false,
		},
		{
			name:  "valid modules",
			path:  "content/modules",
			regex: regexModules,
			want: []*BlueprintMiscContent{
				{
					Name:     "beta-public-cluster",
					Location: "modules/beta-public-cluster",
				},
				{
					Name:     "binary-authorization",
					Location: "modules/binary-authorization",
				},
				{
					Name:     "private-cluster",
					Location: "modules/private-cluster",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid dir",
			path:    "content/modules2",
			regex:   regexModules,
			wantErr: true,
		},
		{
			name:  "some example folders without any tf",
			path:  "content/examples-some-without-tf/examples",
			regex: regexExamples,
			want: []*BlueprintMiscContent{
				{
					Name:     "terraform",
					Location: "examples/acm/acm-terraform-blog-part1/terraform",
				},
				{
					Name:     "simple_regional",
					Location: "examples/simple_regional",
				},
			},
			wantErr: false,
		},
		{
			name:    "all module folders without any tf",
			path:    "content/modules-no-tf/modules",
			regex:   regexModules,
			want:    []*BlueprintMiscContent{},
			wantErr: false,
		},
		{
			name:    "mismatched regex",
			path:    "content/modules",
			regex:   "badRegex",
			want:    []*BlueprintMiscContent{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := regexp.MustCompile(tt.regex)
			sort.SliceStable(tt.want, func(i, j int) bool { return tt.want[i].Name < tt.want[j].Name })
			got, err := getDirPaths(path.Join(bptestdataPath, tt.path), re)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDirPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got, tt.want)
		})
	}
}
