package bpmetadata

import (
	"errors"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/iancoleman/strcase"
)

type repoDetail struct {
	RepoName   string
	ModuleName string
	Source     *repoSource
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
	// For a submodule, we'll try to get repo details from the
	// root blueprint or just return the current repoDetail object
	// if it's still in memory.
	if strings.Contains(bpPath, nestedBpPath) {
		// try to parse the module name from MD which will get
		// overriden with "["repoName-submoduleName" if repoName is available
		r.ModuleName = parseRepoNameFromMd(readme)
		if r.RepoName != "" {
			r.ModuleName = r.RepoName + "-" + getBpSubmoduleName(bpPath, true)
		}

		return
	}

	s := "git"
	bpRootPath := getBlueprintRootPath(bpPath)
	bpPath = strings.TrimSuffix(bpPath, "/")
	repoUrl, repoRoot, err := util.GetRepoUrlAndRootPath(bpPath)
	if err != nil {
		repoUrl = ""
		s = ""
	}

	n, err := util.GetRepoName(repoUrl)
	if err != nil {
		n = parseRepoNameFromMd(readme)
	}

	*r = repoDetail{
		RepoName:   n,
		ModuleName: n,
		Source: &repoSource{
			URL:               repoUrl,
			SourceType:        s,
			BlueprintRootPath: bpRootPath,
			RepoRootPath:      repoRoot,
		},
	}
}

// getRepoDetailsFromRootBp tries to parse repo details from the
// root blueprint metadata.yaml. This is specially useful when
// metadata is generated for a submodule that
func getRepoDetailsFromRootBp(bpPath string) repoDetail {
	rootBp := getBlueprintRootPath(bpPath)
	b, err := UnmarshalMetadata(rootBp, metadataFileName)
	if errors.Is(err, os.ErrNotExist) {
		return repoDetail{
			Source: &repoSource{
				BlueprintRootPath: rootBp,
			},
		}
	}

	// There is metadata for root but does not have source info
	// which means this is a non-git hosted blueprint
	if b.Spec.Info.Source == nil {
		return repoDetail{
			RepoName: b.ResourceMeta.ObjectMeta.NameMeta.Name,
			Source: &repoSource{
				BlueprintRootPath: rootBp,
			},
		}
	}

	// If we get here, root metadata exists and has git info
	return repoDetail{
		RepoName: b.ResourceMeta.ObjectMeta.NameMeta.Name,
		Source: &repoSource{
			URL:               b.Spec.Info.Source.Repo,
			SourceType:        "git",
			BlueprintRootPath: rootBp,
			RepoRootPath:      strings.Replace(rootBp, b.Spec.Info.Source.Dir, "", 1),
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
func getBpSubmoduleName(bpPath string, makeKebab bool) string {
	if strings.Contains(bpPath, nestedBpPath) {
		i := strings.Index(bpPath, nestedBpPath)
		subModuleName := bpPath[i+9:]

		if makeKebab {
			subModuleName = strcase.ToKebab(subModuleName)
		}

		return subModuleName
	}

	return ""
}
