package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_computeFlakesFromBuilds(t *testing.T) {
	tests := []struct {
		name   string
		builds []*build
		want   map[string]flake
	}{
		{
			name: "single",
			builds: []*build{
				getBuild("foo", "j1", "1", "id1", true),
				getBuild("foo", "j1", "1", "id2", false),
			},
			want: map[string]flake{"1-j1": {
				repo:      "foo",
				commitSHA: "1",
				passes:    map[string]*build{"id1": getBuild("foo", "j1", "1", "id1", true)},
				fails:     map[string]*build{"id2": getBuild("foo", "j1", "1", "id2", false)},
			},
			},
		},
		{
			name: "multiple with no flakes",
			builds: []*build{
				getBuild("foo", "j1", "1", "id1", true),
				getBuild("foo", "j1", "2", "id2", false),
			},
			want: map[string]flake{},
		},
		{
			name: "multiple flakes",
			builds: []*build{
				getBuild("foo", "j1", "1", "id1", true),
				getBuild("foo", "j1", "1", "id2", false),
				getBuild("foo", "j1", "1", "id3", false),
				getBuild("foo", "j1", "2", "id4", true),
				getBuild("foo", "j1", "2", "id5", false),
				getBuild("bar", "j1", "2", "id6", false),
				getBuild("bar", "j2", "2", "id7", false),
			},
			want: map[string]flake{
				"1-j1": {
					repo:      "foo",
					commitSHA: "1",
					passes:    map[string]*build{"id1": getBuild("foo", "j1", "1", "id1", true)},
					fails: map[string]*build{
						"id2": getBuild("foo", "j1", "1", "id2", false),
						"id3": getBuild("foo", "j1", "1", "id3", false),
					},
				},
				"2-j1": {
					repo:      "foo",
					commitSHA: "2",
					passes:    map[string]*build{"id4": getBuild("foo", "j1", "2", "id4", true)},
					fails:     map[string]*build{"id5": getBuild("foo", "j1", "2", "id5", false)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeFlakesFromBuilds(tt.builds)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func getBuild(repoName, jobName, commitSHA string, id string, pass bool) *build {
	status := failedStatus
	if pass {
		status = successStatus
	}
	return &build{repoName: repoName, jobName: jobName, commitSHA: commitSHA, id: id, status: status}
}
