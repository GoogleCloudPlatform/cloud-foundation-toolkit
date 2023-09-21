package bpmetadata

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mdTestdataPath = "../testdata/bpmetadata/md"
)

func TestProcessMarkdownContent(t *testing.T) {
	tests := []struct {
		name       string
		fileName   string
		level      int
		order      int
		title      string
		getContent bool
		want       *mdContent
	}{
		{
			name:       "level 1 heading",
			fileName:   "simple-content.md",
			level:      1,
			order:      1,
			getContent: false,
			want: &mdContent{
				literal: "h1 doc title",
			},
		},
		{
			name:       "level 1 heading order 2",
			fileName:   "simple-content.md",
			level:      1,
			order:      2,
			getContent: false,
			want:       nil,
		},
		{
			name:       "level 2 heading order 2",
			fileName:   "simple-content.md",
			level:      2,
			order:      2,
			getContent: false,
			want: &mdContent{
				literal: "Horizontal Rules",
			},
		},
		{
			name:       "level 1 content",
			fileName:   "simple-content.md",
			level:      1,
			order:      1,
			getContent: true,
			want: &mdContent{
				literal: "some content doc title for h1",
			},
		},
		{
			name:       "level 3 content order 2",
			fileName:   "simple-content.md",
			level:      3,
			order:      2,
			getContent: true,
			want: &mdContent{
				literal: "some more content sub heading for h3",
			},
		},
		{
			name:       "content by head title",
			fileName:   "simple-content.md",
			level:      -1,
			order:      -1,
			title:      "h3 sub sub heading",
			getContent: true,
			want: &mdContent{
				literal: "some content sub heading for h3",
			},
		},
		{
			name:       "Tagline does not exist",
			fileName:   "simple-content.md",
			level:      -1,
			order:      -1,
			title:      "Tagline",
			getContent: true,
			want:       nil,
		},
		{
			name:       "Architecture description exists as diagram content",
			fileName:   "list-content.md",
			level:      -1,
			order:      -1,
			title:      "Architecture",
			getContent: true,
			want: &mdContent{
				listItems: []mdListItem{
					{
						text: "User requests are sent to the front end, which is deployed on two Cloud Run services as containers to support high scalability applications.",
					},
					{
						text: "The request then lands on the middle tier, which is the API layer that provides access to the backend. This is also deployed on Cloud Run for scalability and ease of deployment in multiple languages. This middleware is a Golang based API.",
					},
				},
			},
		},
		{
			name:       "content by head title does not exist",
			fileName:   "simple-content.md",
			level:      -1,
			order:      -1,
			title:      "Horizontal Rules",
			getContent: true,
			want:       nil,
		},
		{
			name:       "content by head title link list items",
			fileName:   "list-content.md",
			level:      -1,
			order:      -1,
			title:      "Documentation",
			getContent: true,
			want: &mdContent{
				listItems: []mdListItem{
					{
						text: "document-01",
						url:  "http://google.com/doc-01",
					},
					{
						text: "document-02",
						url:  "http://google.com/doc-02",
					},
					{
						text: "document-03",
						url:  "http://google.com/doc-03",
					},
					{
						text: "document-04",
						url:  "http://google.com/doc-04",
					},
				},
			},
		},
		{
			name:       "content by head title list items",
			fileName:   "list-content.md",
			level:      -1,
			order:      -1,
			title:      "Diagrams",
			getContent: true,
			want: &mdContent{
				listItems: []mdListItem{
					{
						text: "text-document-01",
					},
					{
						text: "text-document-02",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(path.Join(mdTestdataPath, tt.fileName))
			require.NoError(t, err)
			got, _ := getMdContent(content, tt.level, tt.order, tt.title, tt.getContent)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessArchitectureContent(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		title       string
		want        *BlueprintArchitecture
		wantErr     bool
		wantFileErr bool
	}{
		{
			name:     "Architecture details exists as BlueprintArchitecture",
			fileName: "list-content.md",
			title:    "Architecture",
			want: &BlueprintArchitecture{
				Description: []string{
					`1. Step 1`,
					`2. Step 2`,
					`3. Step 3`,
				},
				DiagramUrl: "https://i.redd.it/w3kr4m2fi3111.png",
			},
		},
		{
			name:     "Architecture details don't exist as BlueprintArchitecture",
			fileName: "list-content.md",
			title:    "ArchitectureNotValid",
			wantErr:  true,
		},
		{
			name:        "md content file path for BlueprintArchitecture is invalid",
			fileName:    "list-content-bad-file-name.md",
			title:       "Architecture",
			wantErr:     true,
			wantFileErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(path.Join(mdTestdataPath, tt.fileName))
			if (err != nil) != tt.wantFileErr {
				t.Errorf("ReadFile() = %v, wantErr %v", err, tt.wantFileErr)
				return
			}

			got, err := getArchitctureInfo(content, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("getArchitctureInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessDeploymentDurationContent(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		title    string
		want     *BlueprintTimeEstimate
		wantErr  bool
	}{
		{
			name:     "Deployment duration details exists as BlueprintTimeEstimate",
			fileName: "simple-content.md",
			title:    "Deployment Duration",
			want: &BlueprintTimeEstimate{
				ConfigurationSecs: 120,
				DeploymentSecs:    600,
			},
		},
		{
			name:     "Deployment duration details don't exist as BlueprintTimeEstimate",
			fileName: "simple-content.md",
			title:    "Deployment Duration Invalid",
			wantErr:  true,
		},
		{
			name:     "Deployment duration exists but only for configuration",
			fileName: "simple-content.md",
			title:    "Deployment Duration Only Config",
			want: &BlueprintTimeEstimate{
				ConfigurationSecs: 120,
			},
		},
		{
			name:     "md content file path for BlueprintTimeEstimate is invalid",
			fileName: "simple-content-bad-file-name.md",
			title:    "Does not matter",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, _ := os.ReadFile(path.Join(mdTestdataPath, tt.fileName))
			got, err := getDeploymentDuration(content, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDeploymentDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessCostEstimateContent(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		title    string
		want     *BlueprintCostEstimate
		wantErr  bool
	}{
		{
			name:     "Cost estimate details exists as BlueprintCostEstimate",
			fileName: "simple-content.md",
			title:    "Cost",
			want: &BlueprintCostEstimate{
				Description: "Solution cost details",
				Url:         "https://cloud.google.com/products/calculator?id=02fb0c45-cc29-4567-8cc6-f72ac9024add",
			},
		},
		{
			name:     "Cost estimate details don't exist as BlueprintCostEstimate",
			fileName: "simple-content.md",
			title:    "Cost Invalid",
			wantErr:  true,
		},
		{
			name:     "md content file path for BlueprintCostEstimate is invalid",
			fileName: "simple-content-bad-file-name.md",
			title:    "Does not matter",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, _ := os.ReadFile(path.Join(mdTestdataPath, tt.fileName))
			got, err := getCostEstimate(content, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCostEstimate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
