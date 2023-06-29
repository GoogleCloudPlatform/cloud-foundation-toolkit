package util

import (
	"os"
	"path"
	"strings"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func TestGetRepoUrlAndRootPath(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		subDir  string
		remote  string
		wantURL string
		wantErr bool
	}{
		{
			name:    "simple",
			repo:    "https://github.com/foo/bar",
			remote:  defaultRemote,
			wantURL: "https://github.com/foo/bar",
		},
		{
			name:    "simple trailing",
			repo:    "https://gitlab.com/foo/bar/",
			remote:  defaultRemote,
			wantURL: "https://gitlab.com/foo/bar/",
		},
		{
			name:    "no scheme",
			repo:    "github.com/foo/bar",
			remote:  defaultRemote,
			wantURL: "github.com/foo/bar",
		},
		{
			name:    "invalid remote",
			repo:    "github.com/foo/bar",
			remote:  "foo",
			wantErr: true,
		},
		{
			name:    "simple w/ module sub directory",
			repo:    "https://github.com/foo/bar",
			subDir:  "modules/bp1",
			remote:  defaultRemote,
			wantURL: "https://github.com/foo/bar",
		},
		{
			name:    "simple w/ ssh remote",
			repo:    "git@github.com:foo/bar.git",
			remote:  defaultRemote,
			wantURL: "https://github.com/foo/bar.git",
		},
		{
			name:    "non git@github.com ssh remote",
			repo:    "git@githubAcom:foo/bar.git",
			remote:  defaultRemote,
			wantURL: "git@githubAcom:foo/bar.git",
		},
		{
			name:    "simple w/ module sub directory w/ ssh remote",
			repo:    "git@github.com:foo/bar.git",
			remote:  defaultRemote,
			subDir:  "modules/bp1",
			wantURL: "https://github.com/foo/bar.git",
		},
		{
			name:    "gitlab repo url should not be modified",
			repo:    "git@gitlab.com:foo/bar.git",
			remote:  defaultRemote,
			wantURL: "git@gitlab.com:foo/bar.git",
		},
		{
			name:    "empty repo url",
			repo:    "",
			remote:  defaultRemote,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tempGitRepoWithRemote(t, tt.repo, tt.remote, tt.subDir)
			gotURL, gotPath, err := GetRepoUrlAndRootPath(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepoUrlAndRootPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotURL != tt.wantURL {
				t.Errorf("URL - GetRepoUrlAndRootPath() = %v, want %v", gotURL, tt.wantURL)
			}

			wantPath := strings.TrimSuffix(strings.ReplaceAll(dir, tt.subDir, ""), "/")
			if tt.wantErr {
				wantPath = ""
			}

			if gotPath != wantPath {
				t.Errorf("RootPath - GetRepoUrlAndRootPath() = %v, want %v", gotPath, wantPath)
			}
		})
	}
}

func TestGetRepoNameFromUrl(t *testing.T) {
	tests := []struct {
		name    string
		repoUrl string
		want    string
		wantErr bool
	}{
		{
			name:    "simple",
			repoUrl: "https://github.com/foo/bar",
			want:    "bar",
		},
		{
			name:    "no scheme",
			repoUrl: "github.com/foo/bar",
			want:    "bar",
		},
		{
			name:    "gerrit repo",
			repoUrl: "sso://team/foo/bar",
			want:    "bar",
		},
		{
			name:    "empty Url",
			repoUrl: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRepoName(tt.repoUrl)
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

func tempGitRepoWithRemote(t *testing.T, repoURL, remote string, subDir string) string {
	t.Helper()
	dir := t.TempDir()
	if subDir != "" {
		err := os.MkdirAll(path.Join(dir, subDir), 0755)
		if err != nil {
			t.Fatalf("Error sub dir for temp git repo: %v", err)
		}
	}

	r, err := git.PlainInit(dir, false)
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

	return path.Join(dir, subDir)
}
