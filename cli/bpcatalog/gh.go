package bpcatalog

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

const ghTokenEnvVar = "GITHUB_TOKEN"

type ghService struct {
	client *github.Client
	ctx    context.Context
	orgs   []string
}

type ghServiceOption func(*ghService)

func withOrgs(orgs []string) ghServiceOption {
	return func(g *ghService) {
		g.orgs = orgs
	}
}

func withClient(c *http.Client) ghServiceOption {
	return func(g *ghService) {
		g.client = github.NewClient(c)
	}
}

func withTokenClient() ghServiceOption {
	return func(g *ghService) {
		pat, isSet := os.LookupEnv(ghTokenEnvVar)
		if !isSet {
			Log.Crit(fmt.Sprintf("GitHub token env var %s is not set", ghTokenEnvVar))
			os.Exit(1)
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: pat},
		)
		tc := oauth2.NewClient(g.ctx, ts)
		g.client = github.NewClient(tc)
	}
}

func newGHService(opts ...ghServiceOption) *ghService {
	ctx := context.Background()
	ghs := &ghService{
		client: github.NewClient(nil),
		ctx:    ctx,
	}
	for _, opt := range opts {
		opt(ghs)
	}
	return ghs
}

type repos []*github.Repository

// filter filters repos using a given filter func.
func (r repos) filter(filter func(*github.Repository) bool) repos {
	var filtered []*github.Repository
	for _, repo := range r {
		if filter(repo) {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

// sort sorts repos using a given sort option.
func (r repos) sort(s sortOption) (repos, error) {
	switch s {
	case sortCreated:
		sort.SliceStable(r, func(i, j int) bool { return r[i].GetCreatedAt().Before(r[j].GetCreatedAt().Time) })
	case sortStars:
		sort.SliceStable(r, func(i, j int) bool { return r[i].GetStargazersCount() < r[j].GetStargazersCount() })
	case sortName:
		sort.SliceStable(r, func(i, j int) bool { return r[i].GetName() < r[j].GetName() })
	default:
		return nil, fmt.Errorf("one of %+v expected. unknown format: %s", sortOptions, catalogListFlags.sort)
	}
	return r, nil
}

// fetchRepos fetches all repos across multiple orgs.
func (g *ghService) fetchRepos() (repos, error) {
	opts := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 100}, Type: "public"}
	var allRepos []*github.Repository
	for _, org := range g.orgs {
		for {
			repos, resp, err := g.client.Repositories.ListByOrg(g.ctx, org, opts)
			if err != nil {
				return nil, err
			}
			allRepos = append(allRepos, repos...)
			// if no next page, we have reached end of pagination
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
	}
	return allRepos, nil
}
