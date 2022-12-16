package util

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
)

const defaultRemote = "origin"

// getRepoName finds upstream repo name from a given repo directory
func GetRepoName(dir string) (string, error) {
	remoteUrl, err := GetRepoUrl(dir)
	if err != nil {
		return "", fmt.Errorf("error getting remote URL: %w", err)
	}

	trimmedRemotePath := strings.TrimSuffix(remoteUrl.Path, "/")
	splitRemotePath := strings.Split(trimmedRemotePath, "/")
	// expect path to be /owner/repo
	if len(splitRemotePath) != 3 {
		return "", fmt.Errorf("expected owner/repo, got %s", trimmedRemotePath)
	}

	repoName := strings.TrimSuffix(splitRemotePath[len(splitRemotePath)-1], ".git")
	return repoName, nil
}

// getRepoName finds upstream repo name from a given repo directory
func GetRepoUrl(dir string) (*url.URL, error) {
	opt := &git.PlainOpenOptions{DetectDotGit: true}
	r, err := git.PlainOpenWithOptions(dir, opt)
	if err != nil {
		return nil, fmt.Errorf("error opening git dir %s: %w", dir, err)
	}
	rm, err := r.Remote(defaultRemote)
	if err != nil {
		return nil, fmt.Errorf("error finding remote %s in git dir %s: %w", defaultRemote, dir, err)
	}

	// validate remote URL
	remoteURL, err := url.Parse(rm.Config().URLs[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing remote URL: %w", err)
	}

	return remoteURL, nil
}
