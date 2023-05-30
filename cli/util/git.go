package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	git "github.com/go-git/go-git/v5"
)

const defaultRemote = "origin"

// getRepoName finds upstream repo name from a given repo directory
func GetRepoName(repoUrl string) (string, error) {
	if repoUrl == "" {
		return "", fmt.Errorf("empty URL")
	}

	u, err := url.Parse(repoUrl)
	if err != nil {
		return "", fmt.Errorf("malformed repo URL: %w", err)
	}

	trimmedRemotePath := strings.TrimSuffix(u.Path, "/")
	i := strings.LastIndex(trimmedRemotePath, "/")
	repoName := strings.TrimSuffix(trimmedRemotePath[i+1:], ".git")

	return repoName, nil
}

// GetRepoUrlAndRootPath finds upstream repo URL and the root local path
func GetRepoUrlAndRootPath(dir string) (string, string, error) {
	opt := &git.PlainOpenOptions{DetectDotGit: true}
	r, err := git.PlainOpenWithOptions(dir, opt)
	if err != nil {
		return "", "", fmt.Errorf("error opening git dir %s: %w", dir, err)
	}

	repoRootPath := ""
	repoURL := ""
	rm, err := r.Remote(defaultRemote)
	if err != nil {
		return repoURL, repoRootPath, fmt.Errorf("error finding remote %s in git dir %s: %w", defaultRemote, dir, err)
	}

	if len(rm.Config().URLs) > 0 {
		repoURL = resolveRemoteSSHURLToHTTPS(rm.Config().URLs[0])
	}

	if repoURL == "" {
		return repoURL, repoRootPath, fmt.Errorf("empty URL")
	}

	w, _ := r.Worktree()
	repoRootPath = w.Filesystem.Root()

	// validate remote URL
	_, err = url.Parse(repoURL)
	if err != nil {
		return repoURL, repoRootPath, fmt.Errorf("error parsing remote URL: %w", err)
	}

	return repoURL, repoRootPath, nil
}

func resolveRemoteSSHURLToHTTPS(URL string) string {
	githubSSHRegex := regexp.MustCompile(`git@github.com:`)
	resolvedURL := githubSSHRegex.ReplaceAllString(URL, "https://github.com/")
	return strings.TrimSuffix(resolvedURL, ".git")
}
