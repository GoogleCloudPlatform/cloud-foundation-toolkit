package bpmetadata

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

type mdContent struct {
	literal   string
	listItems []mdListItem
}

type mdListItem struct {
	text string
	url  string
}

// getMdContent accepts 3 types of content requests and return and mdContent object
// with the relevant content info. The 3 scenarios are:
// 1: get heading literal by (level and/or order) OR by title
// 2: get paragraph content immediately following a heading by (level and/or order) OR by title
// 3: get list item content immediately following a heading by (level and/or order) OR by title
// A -1 value to headLevel/headOrder enforces the content to be matchd by headTitle
func getMdContent(content []byte, headLevel int, headOrder int, headTitle string, getContent bool) *mdContent {
	mdDocument := markdown.Parse(content, nil)

	if mdDocument == nil {
		return nil
	}

	orderCtr := 0
	mdSections := mdDocument.GetChildren()
	var foundHead bool
	for _, section := range mdSections {
		currLeaf := ast.GetFirstChild(section).AsLeaf()
		switch sectionType := section.(type) {
		case *ast.Heading:
			if headTitle == string(currLeaf.Literal) {
				foundHead = true
			}

			if headLevel == sectionType.Level {
				orderCtr++
			}

			if !getContent && (headOrder == orderCtr || foundHead) {
				return &mdContent{
					literal: string(currLeaf.Literal),
				}
			}

		case *ast.Paragraph:
			if getContent && (headOrder == orderCtr || foundHead) {
				return &mdContent{
					literal: string(currLeaf.Literal),
				}
			}

		case *ast.List:
			if getContent && (headOrder == orderCtr || foundHead) {
				var mdListItems []mdListItem
				for _, c := range sectionType.Container.Children {
					var listItem mdListItem
					// each item is a list with data and metadata about the list item
					itemConfigs := ast.GetFirstChild(c).AsContainer().Children
					// if the length of the child node is 1, it is a plain text list item
					// if the length is greater the 1, it is a list item with a link
					if len(itemConfigs) == 1 {
						listItemText := string(itemConfigs[0].AsLeaf().Literal)
						listItem = mdListItem{
							text: listItemText,
						}
					} else if len(itemConfigs) > 1 {
						// the second child node has the link data and metadata
						listItemLink := itemConfigs[1].(*ast.Link)
						listItemText := string(ast.GetFirstChild(listItemLink).AsLeaf().Literal)

						listItem = mdListItem{
							text: listItemText,
							url:  string(listItemLink.Destination),
						}
					}

					mdListItems = append(mdListItems, listItem)
				}

				return &mdContent{
					listItems: mdListItems,
				}
			}
		}
	}

	return nil
}
