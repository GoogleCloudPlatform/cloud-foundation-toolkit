package bpbuild

import (
	"fmt"
	"net/url"
	"strings"

	git "github.com/go-git/go-git/v5"
)

const defaultRemote = "origin"

// getRepoName finds upstream repo name from a given repo directory
func getRepoName(dir string) (string, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return "", fmt.Errorf("error opening git dir %s: %w", dir, err)
	}
	rm, err := r.Remote(defaultRemote)
	if err != nil {
		return "", fmt.Errorf("error finding remote %s in git dir %s: %w", defaultRemote, dir, err)
	}

	// validate remote URL
	remoteURL, err := url.Parse(rm.Config().URLs[0])
	if err != nil {
		return "", fmt.Errorf("error parsing remote URL: %w", err)
	}
	trimmedRemotePath := strings.TrimSuffix(remoteURL.Path, "/")
	splitRemotePath := strings.Split(trimmedRemotePath, "/")
	// expect path to be /owner/repo
	if len(splitRemotePath) != 3 {
		return "", fmt.Errorf("expected owner/repo, got %s", trimmedRemotePath)
	}
	return splitRemotePath[len(splitRemotePath)-1], nil
}
