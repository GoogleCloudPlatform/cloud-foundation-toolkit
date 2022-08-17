package bpbuild

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
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
	repoOwnerName := strings.TrimSuffix(remoteURL.Path, "/")
	splitRepoOwnerName := strings.Split(repoOwnerName, "/")
	// expect path to be /owner/repo
	if len(splitRepoOwnerName) != 3 {
		return "", fmt.Errorf("expected owner/repo, got %s", repoOwnerName)
	}
	return splitRepoOwnerName[len(splitRepoOwnerName)-1], nil
}
