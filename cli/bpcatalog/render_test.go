package bpcatalog

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v47/github"
	"github.com/stretchr/testify/assert"
)

const (
	expectedSuffix = ".expected"
	updateEnvVar   = "UPDATE_EXPECTED"
	testDataDir    = "../testdata/catalog"
)

func TestRender(t *testing.T) {
	testRepoData := repos{
		{
			Name:            github.String("terraform-google-bar"),
			CreatedAt:       &github.Timestamp{Time: time.Date(2021, 1, 3, 4, 3, 0, 0, time.UTC)},
			StargazersCount: github.Int(5),
		},
		{
			Name:            github.String("terraform-google-foo"),
			CreatedAt:       &github.Timestamp{Time: time.Date(2022, 11, 3, 4, 3, 0, 0, time.UTC)},
			StargazersCount: github.Int(10),
		},
	}
	tests := []struct {
		name    string
		r       repos
		format  renderFormat
		wantErr bool
	}{
		{
			name:   "table",
			r:      testRepoData,
			format: renderTable,
		},
		{
			name:   "csv",
			r:      testRepoData,
			format: renderCSV,
		},
		{
			name:    "invalid",
			r:       testRepoData,
			format:  "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer
			if err := render(tt.r, &got, tt.format); (err != nil) != tt.wantErr {
				t.Errorf("render() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				expectedPath := path.Join(testDataDir, tt.name+".expected")
				expected := readFile(t, expectedPath)
				updateExpected(t, expectedPath, got.String())
				assert.Equal(t, expected, got.String())
			}
		})
	}
}

func readFile(t *testing.T, p string) string {
	t.Helper()
	j, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("error reading file %s: %s", p, err)
	}
	return string(j)
}

// UpdateExpected updates expected file at fp with data with update env var is set.
func updateExpected(t *testing.T, fp, data string) {
	t.Helper()
	if strings.ToLower(os.Getenv(updateEnvVar)) != "true" {
		return
	}
	// 0755 allows read/execute for everyone, write for owner
	// which is a safe default since this is test data.
	// Execute bit is needed to traverse directories.
	err := os.MkdirAll(path.Dir(fp), 0755)
	if err != nil {
		t.Fatalf("error updating result: %v", err)
	}
	err = os.WriteFile(fp, []byte(data), 0755)
	if err != nil {
		t.Fatalf("error updating result: %v", err)
	}
}
