package bpcatalog

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v68/github"
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
			Name:            github.Ptr("terraform-google-bar"),
			CreatedAt:       &github.Timestamp{Time: time.Date(2021, 1, 3, 4, 3, 0, 0, time.UTC)},
			StargazersCount: github.Ptr(5),
			Description:     github.Ptr("lorem ipsom"),
		},
		{
			Name:            github.Ptr("terraform-google-foo"),
			CreatedAt:       &github.Timestamp{Time: time.Date(2022, 11, 3, 4, 3, 0, 0, time.UTC)},
			StargazersCount: github.Ptr(10),
		},
		{
			Name:            github.Ptr("terraform-foo"),
			CreatedAt:       &github.Timestamp{Time: time.Date(2022, 11, 3, 4, 3, 0, 0, time.UTC)},
			StargazersCount: github.Ptr(10),
			Topics:          []string{"unrelated", e2eLabel, "containers"},
		},
	}
	tests := []struct {
		name    string
		r       repos
		format  renderFormat
		verbose bool
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
			name:    "csv-verbose",
			r:       testRepoData,
			format:  renderCSV,
			verbose: true,
		},
		{
			name:   "html",
			r:      testRepoData,
			format: renderHTML,
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
			if err := render(tt.r, &got, tt.format, tt.verbose); (err != nil) != tt.wantErr {
				t.Errorf("render() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				expectedPath := path.Join(testDataDir, tt.name+".expected")
				updateExpected(t, expectedPath, got.String())
				expected := readFile(t, expectedPath)
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

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		_, err := os.Create(fp)
		if err != nil {
			t.Fatalf("error creating %s: %v", fp, err)
		}
	}

	err = os.WriteFile(fp, []byte(data), 0755)
	if err != nil {
		t.Fatalf("error updating result: %v", err)
	}
}

func TestDocSort(t *testing.T) {
	tests := []struct {
		name  string
		input []displayMeta
		want  []displayMeta
	}{
		{
			name: "simple",
			input: []displayMeta{
				{
					DisplayName: "a",
					IsE2E:       false,
				},
				{
					DisplayName: "b",
					IsE2E:       true,
				},
			},
			want: []displayMeta{
				{
					DisplayName: "b",
					IsE2E:       true,
				},
				{
					DisplayName: "a",
					IsE2E:       false,
				},
			},
		},

		{
			name: "mutiple",
			input: []displayMeta{
				{
					DisplayName: "d",
					IsE2E:       true,
				},
				{
					DisplayName: "b",
					IsE2E:       false,
				},
				{
					DisplayName: "c",
					IsE2E:       false,
				},
				{
					DisplayName: "a",
					IsE2E:       true,
				},
			},
			want: []displayMeta{
				{
					DisplayName: "a",
					IsE2E:       true,
				},
				{
					DisplayName: "d",
					IsE2E:       true,
				},
				{
					DisplayName: "b",
					IsE2E:       false,
				},
				{
					DisplayName: "c",
					IsE2E:       false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := docSort(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReposToDisplayMeta(t *testing.T) {
	tests := []struct {
		name  string
		input repos
		want  []displayMeta
	}{
		{
			name: "simple",
			input: repos{
				{
					Name:            github.Ptr("terraform-google-bar"),
					CreatedAt:       &github.Timestamp{Time: time.Date(2021, 1, 3, 4, 3, 0, 0, time.UTC)},
					StargazersCount: github.Ptr(5),
					Description:     github.Ptr("lorem ipsom"),
					Topics:          []string{"containers"},
				},
				{
					Name:            github.Ptr("terraform-foo"),
					CreatedAt:       &github.Timestamp{Time: time.Date(2022, 11, 3, 4, 3, 0, 0, time.UTC)},
					StargazersCount: github.Ptr(10),
					Topics:          []string{"unrelated", e2eLabel, "containers"},
				},
				{
					Name:            github.Ptr("foo"),
					CreatedAt:       &github.Timestamp{Time: time.Date(2022, 11, 3, 4, 3, 0, 0, time.UTC)},
					StargazersCount: github.Ptr(10),
				},
			},
			want: []displayMeta{
				{
					Name:        "terraform-google-bar",
					DisplayName: "bar",
					Stars:       "5",
					CreatedAt:   "2021-01-03",
					Description: "lorem ipsom",
					Labels:      []string{"containers"},
					URL:         "",
					Categories:  "Containers",
					IsE2E:       false,
				},
				{
					Name:        "terraform-foo",
					DisplayName: "foo",
					Stars:       "10",
					CreatedAt:   "2022-11-03",
					Description: "",
					Labels:      []string{"unrelated", e2eLabel, "containers"},
					URL:         "",
					Categories:  "Containers, End-to-end",
					IsE2E:       true,
				},
				{
					Name:        "foo",
					DisplayName: "foo",
					Stars:       "10",
					CreatedAt:   "2022-11-03",
					Description: "",
					Labels:      []string(nil),
					URL:         "",
					Categories:  "",
					IsE2E:       false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reposToDisplayMeta(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderDocHTML(t *testing.T) {
	tests := []struct {
		name  string
		input []displayMeta
	}{
		{
			name: "single-html",
			input: []displayMeta{
				{
					Name:        "terraform-google-bar",
					DisplayName: "bar",
					Stars:       "5",
					CreatedAt:   "2021-01-03",
					Description: "lorem ipsom",
					Labels:      []string{"containers"},
					URL:         "",
					Categories:  "Containers",
					IsE2E:       false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderDocHTML(tt.input)
			expectedPath := path.Join(testDataDir, tt.name+".expected")
			updateExpected(t, expectedPath, got)
			expected := readFile(t, expectedPath)
			assert.Equal(t, expected, got)
		})
	}
}
