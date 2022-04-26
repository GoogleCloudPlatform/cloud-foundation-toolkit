package cmd

import (
	"context"
	"fmt"
	"time"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

const (
	successStatus = "SUCCESS"
	failedStatus  = "FAILURE"
)

// getCBBuildsWithFilter returns a list of cloudbuild builds in projectID within start and end time
func getCBBuildsWithFilter(start, end time.Time, projectID string) ([]*build, error) {
	ctx := context.Background()
	cloudbuildService, err := cloudbuild.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating cloudbuild service: %v", err)
	}

	filter := fmt.Sprintf("create_time>=\"%s\" AND create_time<\"%s\"", formatTimeCB(start), formatTimeCB(end))
	c, err := cloudbuildService.Projects.Builds.List(projectID).Filter(filter).Do()
	if err != nil {
		return nil, fmt.Errorf("error listing builds with filter %s in project %s: %v", filter, projectID, err)
	}
	cbBuilds := c.Builds
	if len(cbBuilds) < 1 {
		return nil, fmt.Errorf("no builds found with filter %s in project %s", filter, projectID)
	}

	for {
		c, err = cloudbuildService.Projects.Builds.List(projectID).Filter(filter).PageToken(c.NextPageToken).Do()
		if err != nil {
			return nil, fmt.Errorf("error retriving next page with token %s: %v", c.NextPageToken, err)
		}
		cbBuilds = append(cbBuilds, c.Builds...)
		if c.NextPageToken == "" {
			break
		}
	}

	builds := []*build{}
	for _, b := range cbBuilds {
		// filter out builds not triggered from source repos
		commit, commitExists := b.Substitutions["COMMIT_SHA"]
		if !commitExists {
			continue
		}
		repoName, repoNameExists := b.Substitutions["REPO_NAME"]
		if !repoNameExists {
			continue
		}
		triggerName, triggerNameExists := b.Substitutions["TRIGGER_NAME"]
		if !triggerNameExists {
			continue
		}
		builds = append(builds, &build{commitSHA: commit, repoName: repoName, jobName: triggerName, id: b.Id, status: b.Status})
	}
	return builds, nil
}
