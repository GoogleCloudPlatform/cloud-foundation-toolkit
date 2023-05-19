package bpmetadata

import (
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/iancoleman/strcase"
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
func getRepoDetailsByPath(bpPath string, sourceUrl *BlueprintRepoDetail, repoName string, readmeContent []byte) *repoDetail {
	rootRepoPath := getBpRootPath(bpPath)
	if sourceUrl == nil {
		bpPath = strings.TrimSuffix(bpPath, "/")
		repoUrl, err := util.GetRepoUrl(bpPath)
		if err != nil {
			repoUrl = ""
		}

		sourceUrl = &BlueprintRepoDetail{
			Repo: repoUrl,
		}
	}

	if repoName == "" {
		n, err := util.GetRepoName(sourceUrl.Repo)
		if err != nil {
			// Try to get the repo name from readme instead.
			title, err := getMdContent(readmeContent, 1, 1, "", false)
			if err == nil {
				n = strcase.ToKebab(title.literal)
			}
		}

		repoName = n
	}

	return &repoDetail{
		Name: repoName,
		Source: &repoSource{
			Path:       sourceUrl.Repo,
			SourceType: "git",
			RootPath:   rootRepoPath,
		},
	}
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

// getBpSubmoduleName gets the submodule name from the blueprint path
// if it lives under the /modules directory
func getBpSubmoduleName(bpPath string) string {
	if strings.Contains(bpPath, nestedBpPath) {
		i := strings.Index(bpPath, nestedBpPath)
		return bpPath[i+9:]
	}

	return ""
}
