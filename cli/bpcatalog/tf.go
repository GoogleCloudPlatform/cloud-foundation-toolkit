package bpcatalog

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v68/github"
)

// sortOption defines the set of sort options for catalog.
type sortOption string

func (s *sortOption) String() string {
	return string(*s)
}

func (s *sortOption) Empty() bool {
	return s.String() == ""
}

func (s *sortOption) Set(v string) error {
	f, err := sortOptionFromString(v)
	if err != nil {
		return err
	}
	*s = f
	return nil
}

func sortOptionFromString(s string) (sortOption, error) {
	format := sortOption(s)
	for _, stage := range sortOptions {
		if format == stage {
			return format, nil
		}
	}
	return "", fmt.Errorf("one of %+v expected. unknown sort option: %s", sortOptions, s)
}

func (r *sortOption) Type() string {
	return "sortOption"
}

const (
	sortStars   sortOption = "stars"
	sortCreated sortOption = "created"
	sortName    sortOption = "name"
)

var (
	sortOptions = []sortOption{sortStars, sortCreated, sortName}
)

// fetchSortedTFRepos returns a slice of repos sorted by sortOpt.
func fetchSortedTFRepos(gh *ghService, sortOpt sortOption) (repos, error) {
	repos, err := gh.fetchRepos()
	if err != nil {
		return nil, fmt.Errorf("error fetching repos: %w", err)
	}
	repos = repos.filter(func(r *github.Repository) bool {
		if r.GetArchived() {
			return false
		}
		return repoAllowList[r.GetName()] || (strings.HasPrefix(r.GetName(), "terraform-google") && !repoIgnoreList[r.GetName()])
	})
	return repos.sort(sortOpt)
}
