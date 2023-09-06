package bpbuild

import (
	"context"
	"fmt"
	"os"
	"time"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
	"gopkg.in/yaml.v3"
)

const (
	successStatus = "SUCCESS"
	failedStatus  = "FAILURE"
)

// getCBBuildsWithFilter returns a list of cloudbuild builds in projectID with a given filter.
// Additional client side filters can be specified via cFilters.
// TODO(bharathkkb): move https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/tree/master/infra/utils/fbf into CLI
func getCBBuildsWithFilter(projectID string, filter string, cFilters []clientBuildFilter) ([]*cloudbuild.Build, error) {
	ctx := context.Background()
	cloudbuildService, err := cloudbuild.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating cloudbuild service: %w", err)
	}

	c, err := cloudbuildService.Projects.Builds.List(projectID).Filter(filter).Do()
	if err != nil {
		return nil, fmt.Errorf("error listing builds with filter %s in project %s: %w", filter, projectID, err)
	}

	cbBuilds := []*cloudbuild.Build{}
	appendClientFilteredBuilds := func(builds []*cloudbuild.Build) {
		for _, b := range builds {
			appendBuild := true
			for _, cFilter := range cFilters {
				// skip if any client side filter evaluates to false
				if !cFilter(b) {
					appendBuild = false
					break
				}
			}
			if appendBuild {
				cbBuilds = append(cbBuilds, b)
			}
		}
	}

	if len(c.Builds) < 1 {
		return nil, fmt.Errorf("no builds found with filter %s in project %s", filter, projectID)
	}
	appendClientFilteredBuilds(c.Builds)

	// pagination
	for {
		c, err = cloudbuildService.Projects.Builds.List(projectID).Filter(filter).PageToken(c.NextPageToken).Do()
		if err != nil {
			return nil, fmt.Errorf("error retrieving next page with token %s: %w", c.NextPageToken, err)
		}
		appendClientFilteredBuilds(c.Builds)
		if c.NextPageToken == "" {
			break
		}
	}
	return cbBuilds, nil
}

// clientside filter functions
type clientBuildFilter func(*cloudbuild.Build) bool

// filterRealBuilds filters out builds not triggered from source repos (i.e by automation).
func filterRealBuilds(b *cloudbuild.Build) bool {
	for _, subs := range []string{"COMMIT_SHA", "REPO_NAME", "TRIGGER_NAME"} {
		_, substExists := b.Substitutions[subs]
		if !substExists {
			return false
		}
	}
	return true
}

// filterGHRepoBuilds filters builds from a particular repo name.
// TODO:(bharathkkb): We should ideally be using a sever side filter for this https://cloud.google.com/build/docs/view-build-results#filtering_build_results_using_queries
// but I was not able to figure out expected format for GH URLs.
func filterGHRepoBuilds(repo string) clientBuildFilter {
	return func(b *cloudbuild.Build) bool {
		name, exists := b.Substitutions["REPO_NAME"]
		if !exists {
			return false
		}
		return name == repo
	}
}

// successBuildsBtwFilterExpr returns a CEL expression as string
// for finding all successful builds between start and end time.
func successBuildsBtwFilterExpr(start, end time.Time) string {
	return fmt.Sprintf(
		"create_time>=\"%s\" AND create_time<\"%s\" AND status=\"%s\"",
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
		successStatus)
}

// getBuildFromFile unmarshalls a CloudBuild file at path.
func getBuildFromFile(path string) (*cloudbuild.Build, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var b cloudbuild.Build
	err = yaml.Unmarshal(content, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// getBuildStepIDs retrieves a slice of build step IDs in a build.
func getBuildStepIDs(b *cloudbuild.Build) []string {
	steps := []string{}
	for _, bs := range b.Steps {
		steps = append(steps, bs.Id)
	}
	return steps
}

// findBuildStageDurations computes duration for a given build stage across a slice of builds
// if and only if stage is successful.
func findBuildStageDurations(stepId string, builds []*cloudbuild.Build) ([]time.Duration, error) {
	durations := []time.Duration{}
	for _, b := range builds {
		for _, bs := range b.Steps {
			if bs.Id != stepId || bs.Status != successStatus {
				continue
			}

			parsedStartTime, err := time.Parse(time.RFC3339Nano, bs.Timing.StartTime)
			if err != nil {
				return []time.Duration{}, err
			}
			parsedEndTime, err := time.Parse(time.RFC3339Nano, bs.Timing.EndTime)
			if err != nil {
				return []time.Duration{}, err
			}
			durations = append(durations, parsedEndTime.Sub(parsedStartTime).Truncate(time.Second))
		}
	}
	return durations, nil
}
