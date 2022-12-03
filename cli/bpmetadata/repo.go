package bpmetadata

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
)

type repoDetail struct {
	Name   string
	Source *repoSource
}

type repoSource struct {
	Path       string
	RootPath   string
	SourceType string
}

const (
	nestedBpPath = "/modules"
)

// getRepoDetailsByPath takes a local path for a blueprint and tries
// to get repo details that include its name, path and type
func getRepoDetailsByPath(bpPath string) (*repoDetail, error) {
	bpPath = strings.TrimSuffix(bpPath, "/")
	rootRepoPath := getBpRootPath(bpPath)
	repoName, err := util.GetRepoName(rootRepoPath)
	if err != nil {
		return nil, fmt.Errorf("error getting the repo name from the provided local repo path: %w", err)
	}

	repoUrl, err := util.GetRepoUrl(bpPath)
	if err != nil {
		return nil, fmt.Errorf("error getting the repo URL from the provided local repo path: %w", err)
	}

	return &repoDetail{
		Name: repoName,
		Source: &repoSource{
			Path:       repoUrl.String(),
			SourceType: "git",
			RootPath:   rootRepoPath,
		},
	}, nil
}

// getBpRootPath determines if the provided bpPath is for a submodule
// and resolves it to the root module path if necessary
func getBpRootPath(bpPath string) string {

	if strings.Contains(bpPath, nestedBpPath) {
		i := strings.Index(bpPath, nestedBpPath)
		bpPath = bpPath[0:i]
	}

	return bpPath
}
