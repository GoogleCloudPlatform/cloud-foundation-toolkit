package bpmetadata

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
)

type repoDetail struct {
	Name   string
	Source *repoSource
}

type repoSource struct {
	Path       string
	SourceType string
}

const (
	bpUrlRegEx   = "(.*terraform-google-[^/]+)/?(modules/.*)?"
	nestedBpPath = "modules/"
)

// getRepoDetailsByPath takes a local path for a blueprint and tries
// to get repo details that include its name, path and type
func getRepoDetailsByPath(bpPath string) (*repoDetail, error) {
	bpPath = strings.TrimSuffix(bpPath, "/")
	repoPath, err := getBpPathForRepoName(bpPath)
	if err != nil {
		return nil, fmt.Errorf("error getting the repo path from the provided blueprint path: %w", err)
	}

	repoName, err := util.GetRepoName(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error getting the repo name from the provided local repo path: %w", err)
	}

	repoUrl, err := util.GetRepoUrl(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error getting the repo URL from the provided local repo path: %w", err)
	}

	return &repoDetail{
		Name: repoName,
		Source: &repoSource{
			Path:       repoUrl,
			SourceType: "git",
		},
	}, nil
}

// getBpPathForRepoName verifies if the blueprint follows blueprint
// naming conventions and returns the local path for the repo root
func getBpPathForRepoName(bpPath string) (string, error) {
	r := regexp.MustCompile(bpUrlRegEx)
	matches := r.FindStringSubmatch(bpPath)

	// not a valid blueprint path if there is no match
	if matches == nil {
		return "", fmt.Errorf("provided blueprint path is not valid: %s", bpPath)
	}

	// if matched, matches should haveexactly 3 items,
	// [0] for the match and [1] for root repo path and [2]
	// for the nested blueprint name
	if len(matches) != 3 {
		return "", fmt.Errorf("provided nested blueprint path is not valid: %s. It should be under the %s directory", bpPath, nestedBpPath)
	}

	// check if the path has a nested blueprint
	// and under the right directory
	if len(bpPath) != len(matches[1]) && matches[2] == "" {
		return "", fmt.Errorf("provided nested blueprint path is not valid: %s. It should be under the %s directory", bpPath, nestedBpPath)
	}

	return matches[1], nil
}
