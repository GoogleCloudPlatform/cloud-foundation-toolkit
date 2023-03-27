package bpmetadata

import (
	"fmt"
	"strings"

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
func getMdContent(content []byte, headLevel int, headOrder int, headTitle string, getContent bool) (*mdContent, error) {
	mdDocument := markdown.Parse(content, nil)
	orderCtr := 0
	mdSections := mdDocument.GetChildren()
	var foundHead bool
	for _, section := range mdSections {
		// if the first child is nil, it's a comment and we don't
		// need to evaluate it
		if ast.GetFirstChild(section) == nil {
			continue
		}

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
				}, nil
			}

		case *ast.Paragraph:
			if getContent && (headOrder == orderCtr || foundHead) {
				return &mdContent{
					literal: string(currLeaf.Literal),
				}, nil
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
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find md content")
}

// getArchitctureInfo parses and builds Architecture details from README.md
func getArchitctureInfo(content []byte, headTitle string) (*BlueprintArchitecture, error) {
	mdDocument := markdown.Parse(content, nil)
	if mdDocument == nil {
		return nil, fmt.Errorf("unable to parse md content")
	}

	children := mdDocument.GetChildren()
	for _, node := range children {
		h, isHeading := node.(*ast.Heading)
		if !isHeading {
			continue
		}

		// check if this is the architecture heading
		hLiteral := string(ast.GetFirstChild(h).AsLeaf().Literal)
		if hLiteral != headTitle {
			continue
		}

		//get architecture details
		infoNode := ast.GetNextNode(h)
		paraNode, isPara := infoNode.(*ast.Paragraph)
		if !isPara {
			continue
		}

		t := ast.GetLastChild(paraNode)
		_, isText := t.(*ast.Text)
		if !isText {
			continue
		}

		d := strings.TrimLeft(string(t.AsLeaf().Literal), "\n")
		dList := strings.Split(d, "\n")
		i := ast.GetPrevNode(t)
		iNode, isImage := i.(*ast.Image)
		if isImage {
			return &BlueprintArchitecture{
				Description: dList,
				DiagramURL:  string(iNode.Destination),
			}, nil
		}

		lNode, isLink := i.(*ast.Link)
		if isLink {
			return &BlueprintArchitecture{
				Description: dList,
				DiagramURL:  string(lNode.Destination),
			}, nil
		}
	}

	return nil, fmt.Errorf("unable to find architecture content")
}
