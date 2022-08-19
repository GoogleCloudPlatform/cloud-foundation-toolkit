package bpbuild

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

func TestFilterRealBuilds(t *testing.T) {
	tests := []struct {
		name  string
		subst map[string]string
		want  bool
	}{
		{
			name:  "fail",
			subst: map[string]string{"foo": "bar"},
			want:  false,
		},
		{
			name:  "partial",
			subst: map[string]string{"REPO_NAME": "bar"},
			want:  false,
		},
		{
			name: "pass",
			subst: map[string]string{
				"REPO_NAME":    "bar",
				"COMMIT_SHA":   "bar",
				"TRIGGER_NAME": "bar",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newTestBuild(tt.subst, nil)
			if got := filterRealBuilds(b); got != tt.want {
				t.Errorf("filterRealBuilds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterGHRepoBuilds(t *testing.T) {
	tests := []struct {
		name  string
		repo  string
		subst map[string]string
		want  bool
	}{
		{
			name:  "fail",
			repo:  "foo",
			subst: map[string]string{"foo": "bar"},
			want:  false,
		},
		{
			name:  "pass",
			repo:  "foo",
			subst: map[string]string{"REPO_NAME": "foo"},
			want:  true,
		},
		{
			name:  "fail different",
			repo:  "bar",
			subst: map[string]string{"REPO_NAME": "foo"},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newTestBuild(tt.subst, nil)
			if got := filterGHRepoBuilds(tt.repo)(b); got != tt.want {
				t.Errorf("filterGHRepoBuilds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindBuildStageDurations(t *testing.T) {
	tests := []struct {
		name    string
		stepId  string
		builds  []*cloudbuild.Build
		want    []time.Duration
		wantErr bool
	}{
		{
			name:   "simple",
			stepId: "foo",
			builds: []*cloudbuild.Build{newTestBuild(nil, []*cloudbuild.BuildStep{newTestBuildStep("foo", time.Hour, successStatus)})},
			want:   []time.Duration{time.Hour},
		},
		{
			name:   "multiple builds",
			stepId: "foo",
			builds: []*cloudbuild.Build{
				newTestBuild(nil, []*cloudbuild.BuildStep{newTestBuildStep("foo", time.Hour, successStatus)}),
				newTestBuild(nil, []*cloudbuild.BuildStep{newTestBuildStep("foo", time.Hour*2, successStatus)}),
				newTestBuild(nil, []*cloudbuild.BuildStep{newTestBuildStep("foo", time.Hour*4, successStatus)}),
			},
			want: []time.Duration{
				time.Hour,
				time.Hour * 2,
				time.Hour * 4,
			},
		},
		{
			name:   "multiple builds multiple steps",
			stepId: "foo",
			builds: []*cloudbuild.Build{
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour, successStatus),
					newTestBuildStep("bar", time.Hour*2, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*2, successStatus),
					newTestBuildStep("bar", time.Hour*8, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*4, successStatus),
					newTestBuildStep("bar", time.Hour, successStatus),
				}),
			},
			want: []time.Duration{
				time.Hour,
				time.Hour * 2,
				time.Hour * 4,
			},
		},
		{
			name:   "multiple builds multiple steps with fails",
			stepId: "foo",
			builds: []*cloudbuild.Build{
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour, failedStatus),
					newTestBuildStep("bar", time.Hour*2, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*2, successStatus),
					newTestBuildStep("bar", time.Hour*8, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*4, failedStatus),
					newTestBuildStep("bar", time.Hour, successStatus),
				}),
			},
			want: []time.Duration{
				time.Hour * 2,
			},
		},
		{
			name:   "empty multiple builds multiple steps but all matched step failed",
			stepId: "foo",
			builds: []*cloudbuild.Build{
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour, failedStatus),
					newTestBuildStep("bar", time.Hour*2, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*2, failedStatus),
					newTestBuildStep("bar", time.Hour*8, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*4, failedStatus),
					newTestBuildStep("bar", time.Hour, successStatus),
				}),
			},
			want: []time.Duration{},
		},
		{
			name:   "empty multiple builds multiple steps no match",
			stepId: "baz",
			builds: []*cloudbuild.Build{
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour, successStatus),
					newTestBuildStep("bar", time.Hour*2, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*2, successStatus),
					newTestBuildStep("bar", time.Hour*8, successStatus),
				}),
				newTestBuild(nil, []*cloudbuild.BuildStep{
					newTestBuildStep("foo", time.Hour*4, successStatus),
					newTestBuildStep("bar", time.Hour, successStatus),
				}),
			},
			want: []time.Duration{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findBuildStageDurations(tt.stepId, tt.builds)
			if (err != nil) != tt.wantErr {
				t.Errorf("findBuildStageDurations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("findBuildStageDurations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newTestBuild(subst map[string]string, bs []*cloudbuild.BuildStep) *cloudbuild.Build {
	return &cloudbuild.Build{
		Substitutions: subst,
		Steps:         bs,
	}
}

func newTestBuildStep(id string, length time.Duration, status string) *cloudbuild.BuildStep {
	return &cloudbuild.BuildStep{
		Id:     id,
		Status: status,
		Timing: &cloudbuild.TimeSpan{
			StartTime: time.Now().Format(time.RFC3339Nano),
			EndTime:   time.Now().Add(length).Format(time.RFC3339Nano),
		},
	}
}
