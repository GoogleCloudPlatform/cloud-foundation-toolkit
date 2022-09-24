package bpmetadata

import (
	"os"
	"path"
	"testing"
)

const (
	mdTestdataPath = "../testdata/bpmetadata/md"
)

func TestProcessMarkdownContent(t *testing.T) {
	tests := []struct {
		name       string
		filePath   string
		level      int
		order      int
		title      string
		getContent bool
		want       *mdContent
	}{
		{
			name:       "level 1 heading",
			filePath:   path.Join(mdTestdataPath, "simple-content.md"),
			level:      1,
			order:      1,
			getContent: false,
			want: &mdContent{
				literal: "h1 doc title",
			},
		},
		{
			name:       "level 1 heading order 2",
			filePath:   path.Join(mdTestdataPath, "simple-content.md"),
			level:      1,
			order:      2,
			getContent: false,
			want:       nil,
		},
		{
			name:       "level 2 heading order 2",
			filePath:   path.Join(mdTestdataPath, "simple-content.md"),
			level:      2,
			order:      2,
			getContent: false,
			want: &mdContent{
				literal: "Horizontal Rules",
			},
		},
		{
			name:       "level 1 content",
			filePath:   path.Join(mdTestdataPath, "simple-content.md"),
			level:      1,
			order:      1,
			getContent: true,
			want: &mdContent{
				literal: "some content doc title for h1",
			},
		},
		{
			name:       "level 3 content order 2",
			filePath:   path.Join(mdTestdataPath, "simple-content.md"),
			level:      3,
			order:      2,
			getContent: true,
			want: &mdContent{
				literal: "some more content sub heading for h3",
			},
		},
		{
			name:       "content by head title",
			filePath:   path.Join(mdTestdataPath, "simple-content.md"),
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
			filePath:   path.Join(mdTestdataPath, "simple-content.md"),
			level:      -1,
			order:      -1,
			title:      "Horizontal Rules",
			getContent: true,
			want:       nil,
		},
		{
			name:       "content by head title link list items",
			filePath:   path.Join(mdTestdataPath, "list-content.md"),
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
			filePath:   path.Join(mdTestdataPath, "list-content.md"),
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
			content, _ := os.ReadFile(tt.filePath)
			got := getMdContent(content, tt.level, tt.order, tt.title, tt.getContent)

			if got != nil {
				if got.literal != tt.want.literal {
					t.Errorf("getMdContent() = %v, want %v", got.literal, tt.want.literal)
					return
				}

				if len(got.listItems) != len(tt.want.listItems) {
					t.Errorf("getMdContent() = %v list items, want %v list items", len(got.listItems), len(tt.want.listItems))
					return
				}

				for i := 0; i < len(got.listItems); i++ {
					if (got.listItems[i].text != tt.want.listItems[i].text) ||
						(got.listItems[i].url != tt.want.listItems[i].url) {
						t.Errorf("getMdContent() = %v list item, want %v list item", got.listItems[i], tt.want.listItems[i])
					}
				}
			} else {
				if tt.want != nil {
					t.Errorf("getMdContent() = returned nil when we want %v", tt.want)

				}
			}
		})
	}
}
