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
	URL               string
	BlueprintRootPath string
	RepoRootPath      string
	SourceType        string
}

const (
	nestedBpPath = "/modules"
)

// getRepoDetailsByPath takes a local path for a blueprint and tries
// to get repo details that include its name, path and type
func getRepoDetailsByPath(bpPath string, r *repoDetail, readme []byte) {
	if strings.Contains(bpPath, nestedBpPath) {
		return
	}

	bpRootPath := getBlueprintRootPath(bpPath)
	bpPath = strings.TrimSuffix(bpPath, "/")
	repoUrl, repoRoot, err := util.GetRepoUrlAndRootPath(bpPath)
	if err != nil {
		repoUrl = ""
	}

	n, err := util.GetRepoName(repoUrl)
	if err != nil {
		n = parseRepoNameFromMd(readme)
	}

	*r = repoDetail{
		Name: n,
		Source: &repoSource{
			URL:               repoUrl,
			SourceType:        "git",
			BlueprintRootPath: bpRootPath,
			RepoRootPath:      repoRoot,
		},
	}
}

func parseRepoNameFromMd(readme []byte) string {
	n := ""
	title, err := getMdContent(readme, 1, 1, "", false)
	if err == nil {
		n = strcase.ToKebab(title.literal)
	}

	return n
}

// getBpRootPath determines if the provided bpPath is for a submodule
// and resolves it to the root module path if necessary
func getBlueprintRootPath(bpPath string) string {
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
