package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/table"
)

// FlakyFinder finds flakes between start and end times
type FlakyFinder struct {
	startTime time.Time
	endTime   time.Time
	projectID string
	verbose   bool
	flakes    map[string]flake
}

// flake represents a collection of flaky builds for a given commit
type flake struct {
	repo      string            // repo name
	commitSHA string            // commit SHA
	passes    map[string]*build // builds passed for commit
	fails     map[string]*build // builds passed for commit
}

// build represents a single instance of a job invoken at commitSHA on a source repo
type build struct {
	repoName  string
	jobName   string
	commitSHA string
	id        string
	status    string
}

// todo(bharathkkb): use config
func NewFlakyFinder(start, end, projectID string, verbose bool) (*FlakyFinder, error) {
	startTime, err := getTimeFromStr(start)
	if err != nil {
		return nil, fmt.Errorf("error parsing startime: %v", err)
	}
	endTime, err := getTimeFromStr(end)
	if err != nil {
		return nil, fmt.Errorf("error parsing endTime: %v", err)
	}
	if projectID == "" {
		return nil, fmt.Errorf("error got empty project ID")
	}
	return &FlakyFinder{
		startTime: startTime,
		endTime:   endTime,
		projectID: projectID,
		verbose:   verbose,
	}, nil
}

func (f *FlakyFinder) ComputeFlakes() error {
	// get builds
	s := spinner.New(spinner.CharSets[35], 500*time.Millisecond)
	s.Start()
	// todo(bharathkkb): support other build systems
	builds, err := getCBBuildsWithFilter(f.startTime, f.endTime, f.projectID)
	if err != nil {
		return fmt.Errorf("error getting builds: %v", err)
	}
	s.Stop()
	// compute flakes
	f.flakes = computeFlakesFromBuilds(builds)
	return nil
}

// computeFlakesFromBuilds computes flakes from a slice of builds
// a collection of builds are considered flakey iff at least two builds
// have passed and failed at the same commit in a repo when triggered by the same job
func computeFlakesFromBuilds(builds []*build) map[string]flake {
	flakes := make(map[string]flake)
	for _, b1 := range builds {
		// commit may have multiple builds so the key for flake lookup is
		// computed from commitSHA and job name
		flakeKey := fmt.Sprintf("%s-%s", b1.commitSHA, b1.jobName)
		// skip if flakes with same flakeKey were previously computed
		//todo(bharathkkb): optimize, we can probably remove elems in a flake from build slice
		_, exists := flakes[flakeKey]
		if exists {
			continue
		}
		// store individual build info
		passedBuildsWithCommit := make(map[string]*build)
		failedBuildsWithCommit := make(map[string]*build)
		storeBuildInfo := func(b *build) {
			switch b.status {
			case successStatus:
				passedBuildsWithCommit[b.id] = b
			case failedStatus:
				failedBuildsWithCommit[b.id] = b
			}
		}
		storeBuildInfo(b1)

		for _, b2 := range builds {
			// match other builds with same commit,repo and job
			if b1.commitSHA == b2.commitSHA &&
				b1.repoName == b2.repoName &&
				b1.jobName == b2.jobName {
				storeBuildInfo(b2)
			}
		}

		// At least one pass and one fail for a given commit is necessary to become a flake
		if len(passedBuildsWithCommit) > 0 && len(failedBuildsWithCommit) > 0 {
			flakes[flakeKey] = flake{repo: b1.repoName, commitSHA: b1.commitSHA, passes: passedBuildsWithCommit, fails: failedBuildsWithCommit}
		}
	}
	return flakes
}

// render displays results in a tabular format
func (f *FlakyFinder) Render() {
	// verbose table with build ids
	tableVerbose := table.NewWriter()
	tableVerbose.SetOutputMirror(os.Stdout)
	tableVerbose.AppendHeader(table.Row{"Repo", "Commit", "Pass Build IDs", "Fail Build IDs"})
	// flakes per repo
	repoFlakeCount := make(map[string]int)
	// flake failures per repo
	repoFlakeFailCount := make(map[string]int)
	for _, f := range f.flakes {
		repoFlakeCount[f.repo]++
		repoFlakeFailCount[f.repo] += len(f.fails)
		pass := ""
		for id := range f.passes {
			pass += id + "\n"
		}
		fail := ""
		for id := range f.fails {
			fail += id + "\n"
		}
		tableVerbose.AppendRow(table.Row{f.repo, f.commitSHA, pass, fail})
		tableVerbose.AppendSeparator()
	}
	if f.verbose {
		tableVerbose.Render()
	}

	// overview table with total number of flakes per repo
	tableOverview := table.NewWriter()
	tableOverview.SetOutputMirror(os.Stdout)
	tableOverview.AppendHeader(table.Row{"Repo", "Flakes", "Flake Failures"})
	totalFlakeCount := 0
	totalFlakeFailCount := 0
	for repo, flakeCount := range repoFlakeCount {
		tableOverview.AppendRow(table.Row{repo, flakeCount, repoFlakeFailCount[repo]})
		totalFlakeCount += flakeCount
		totalFlakeFailCount += repoFlakeFailCount[repo]
	}
	tableOverview.AppendFooter(table.Row{"Total", totalFlakeCount, totalFlakeFailCount})
	tableOverview.Render()
}
