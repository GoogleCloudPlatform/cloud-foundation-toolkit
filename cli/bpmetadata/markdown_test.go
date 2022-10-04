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
					mdListItem{
						text: "document-01",
						url:  "http://google.com/doc-01",
					},
					mdListItem{
						text: "document-02",
						url:  "http://google.com/doc-02",
					},
					mdListItem{
						text: "document-03",
						url:  "http://google.com/doc-03",
					},
					mdListItem{
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
					mdListItem{
						text: "text-document-01",
					},
					mdListItem{
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
