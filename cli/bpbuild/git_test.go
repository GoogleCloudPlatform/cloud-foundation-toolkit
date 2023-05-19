package bpbuild

import (
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func TestGetRepoName(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		remote  string
		want    string
		wantErr bool
	}{
		{
			name:   "simple",
			repo:   "https://github.com/foo/bar",
			remote: defaultRemote,
			want:   "bar",
		},
		{
			name:   "simple trailing",
			repo:   "https://gitlab.com/foo/bar/",
			remote: defaultRemote,
			want:   "bar",
		},
		{
			name:   "no scheme",
			repo:   "github.com/foo/bar",
			remote: defaultRemote,
			want:   "bar",
		},
		{
			name:    "invalid path",
			repo:    "github.com/foo/bar/baz",
			remote:  defaultRemote,
			wantErr: true,
		},
		{
			name:    "invalid remote",
			repo:    "github.com/foo/bar",
			remote:  "foo",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tempGitRepoWithRemote(t, tt.repo, tt.remote)
			got, err := getRepoName(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRepoName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func tempGitRepoWithRemote(t *testing.T, repoURL, remote string) string {
	t.Helper()
	dir := t.TempDir()
	r, err := git.PlainInit(dir, true)
	if err != nil {
		t.Fatalf("Error creating git repo in tempdir: %v", err)
	}
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{repoURL},
	})
	if err != nil {
		t.Fatalf("Error creating remote in tempdir repo: %v", err)
	}
	return dir
}
