package bpcatalog

import (
	"testing"
	"time"

	"github.com/google/go-github/v68/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestFetchSortedTFRepos(t *testing.T) {
	mockT := time.Now()
	tests := []struct {
		name    string
		repos   []github.Repository
		sortBy  sortOption
		want    []string
		wantErr bool
	}{
		{
			name: "simple sort created",
			repos: []github.Repository{
				{
					Name:      github.Ptr("terraform-google-bar"),
					CreatedAt: &github.Timestamp{Time: mockT.Add(time.Hour * 3)},
				},
				{
					Name:      github.Ptr("terraform-google-foo"),
					CreatedAt: &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
				},
				{
					Name:      github.Ptr("foo"),
					CreatedAt: &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
				},
			},
			want: []string{
				"terraform-google-foo",
				"terraform-google-bar",
			},
			sortBy: sortCreated,
		},
		{
			name: "simple sort name",
			repos: []github.Repository{
				{
					Name:      github.Ptr("terraform-google-bar"),
					CreatedAt: &github.Timestamp{Time: mockT.Add(time.Hour * 3)},
				},
				{
					Name:      github.Ptr("terraform-google-foo"),
					CreatedAt: &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
				},
				{
					Name:      github.Ptr("foo"),
					CreatedAt: &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
				},
			},
			want: []string{
				"terraform-google-bar",
				"terraform-google-foo",
			},
			sortBy: sortName,
		},
		{
			name: "simple sort stars",
			repos: []github.Repository{
				{
					Name:            github.Ptr("terraform-google-bar"),
					CreatedAt:       &github.Timestamp{Time: mockT.Add(time.Hour * 3)},
					StargazersCount: github.Ptr(5),
				},
				{
					Name:            github.Ptr("terraform-google-foo"),
					CreatedAt:       &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
					StargazersCount: github.Ptr(10),
				},
				{
					Name:            github.Ptr("foo"),
					CreatedAt:       &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
					StargazersCount: github.Ptr(12),
				},
				{
					Name:      github.Ptr("archived"),
					CreatedAt: &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
					Archived:  github.Ptr(true),
				},
			},
			want: []string{
				"terraform-google-bar",
				"terraform-google-foo",
			},
			sortBy: sortStars,
		},
		{
			name: "invalid",
			repos: []github.Repository{
				{
					Name:            github.Ptr("terraform-google-bar"),
					CreatedAt:       &github.Timestamp{Time: mockT.Add(time.Hour * 3)},
					StargazersCount: github.Ptr(5),
				},
				{
					Name:            github.Ptr("terraform-google-foo"),
					CreatedAt:       &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
					StargazersCount: github.Ptr(10),
				},
				{
					Name:            github.Ptr("foo"),
					CreatedAt:       &github.Timestamp{Time: mockT.Add(time.Hour * 2)},
					StargazersCount: github.Ptr(12),
				},
			},
			wantErr: true,
			sortBy:  "baz",
		},
		{
			name:    "empty",
			repos:   []github.Repository{},
			wantErr: false,
			sortBy:  "name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetOrgsReposByOrg,
					tt.repos,
				),
			)
			mockGHService := newGHService(withClient(mockedHTTPClient), withOrgs([]string{"foo"}))
			got, err := fetchSortedTFRepos(mockGHService, tt.sortBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchSortedTFRepos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotRepoNames []string
			for _, r := range got {
				gotRepoNames = append(gotRepoNames, r.GetName())
			}
			assert.Equal(t, tt.want, gotRepoNames)
		})
	}
}
