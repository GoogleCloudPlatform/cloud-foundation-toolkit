package bpcatalog

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

// renderFormat defines the set of render options for catalog.
type renderFormat string

func (r *renderFormat) String() string {
	return string(*r)
}

func (r *renderFormat) Empty() bool {
	return r.String() == ""
}

func (r *renderFormat) Set(v string) error {
	f, err := renderFormatFromString(v)
	if err != nil {
		return err
	}
	*r = f
	return nil
}

func renderFormatFromString(s string) (renderFormat, error) {
	format := renderFormat(s)
	for _, stage := range renderFormats {
		if format == stage {
			return format, nil
		}
	}
	return "", fmt.Errorf("one of %+v expected. unknown format: %s", renderFormats, s)
}

func (r *renderFormat) Type() string {
	return "renderFormat"
}

const (
	renderTable renderFormat = "table"
	renderCSV   renderFormat = "csv"
	renderHTML  renderFormat = "html"

	renderTimeformat = "2006-01-02"
	e2eLabel         = "end-to-end"
	htmlTemplate     = `<table>
<thead>
    <tr>
      <th>Category</th>
      <th>Blueprint</th>
      <th>Description</th>
    </tr>
  </thead>
<tbody class="list">
{{range .}}{{if .Categories}}
<tr>
      <td>{{.Categories}}</td>
      <td><a
href="{{.URL}}" class="external">{{.DisplayName}}</a></td>
      <td>{{.Description}}</td>
</tr>
{{end}}{{end}}
</tbody>
</table>`
)

var (
	renderFormats = []renderFormat{renderTable, renderCSV, renderHTML}

	// maps GH topics to categories
	// https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/29e980be9f3e3535f4b0b7314c9e1aea5ec2001f/infra/terraform/test-org/org/locals.tf#L39-L53
	topicToCategory = map[string]string{
		e2eLabel:                   "End-to-end",
		"healthcare-life-sciences": "Healthcare and life sciences",
		"serverless-computing":     "Serverless computing",
		"compute":                  "Compute",
		"containers":               "Containers",
		"databases":                "Databases",
		"networking":               "Networking",
		"data-analytics":           "Data analytics",
		"storage":                  "Storage",
		"operations":               "Operations",
		"developer-tools":          "Developer tools",
		"security-identity":        "Security and identity",
		"workspace":                "Workspace",
	}

	// static display data for docs mode
	// these repos are not currently auto discovered
	staticDM = []displayMeta{
		{
			DisplayName: "fabric",
			URL:         "https://github.com/terraform-google-modules/cloud-foundation-fabric",
			Categories:  "End to end",
			IsE2E:       true,
			Description: "Advanced examples designed for prototyping",
		},
		{
			DisplayName: "ai-notebook",
			URL:         "https://github.com/GoogleCloudPlatform/notebooks-blueprint-security",
			Categories:  "End to end, Data analytics",
			IsE2E:       true,
			Description: "Protect confidential data in Vertex AI Workbench notebooks",
		},
	}
)

// displayMeta stores processed display metadata.
// Currently it processes from repo info but
// may also pull from other sources like blueprint meta
// in the future.
type displayMeta struct {
	Name        string
	DisplayName string
	Stars       string
	CreatedAt   string
	Description string
	Labels      []string
	URL         string
	Categories  string
	IsE2E       bool
}

// render writes given repo information in the specified renderFormat to w.
func render(r repos, w io.Writer, format renderFormat, verbose bool) error {
	dm := reposToDisplayMeta(r)
	if format == renderHTML {
		_, err := w.Write([]byte(renderDocHTML(append(dm, staticDM...))))
		if err != nil {
			return err
		}
		return nil
	}

	tbl := table.NewWriter()
	tbl.SetOutputMirror(w)
	h := table.Row{"Repo", "Stars", "Created"}
	if verbose {
		h = append(h, "Description")
	}
	tbl.AppendHeader(h)

	for _, repo := range r {
		row := table.Row{repo.GetName(), repo.GetStargazersCount(), repo.GetCreatedAt().Format(renderTimeformat)}
		if verbose {
			row = append(row, repo.GetDescription())
		}
		tbl.AppendRow(row)
	}
	switch format {
	case renderTable:
		tbl.Render()
	case renderCSV:
		tbl.RenderCSV()
	default:
		return fmt.Errorf("one of %+v expected. unknown format: %s", renderFormats, catalogListFlags.format)
	}
	return nil
}

// reposToDisplayMeta converts repo to displayMeta.
func reposToDisplayMeta(r repos) []displayMeta {
	dm := make([]displayMeta, 0, len(r))
	for _, repo := range r {
		displayName := strings.TrimPrefix(repo.GetName(), "terraform-google-")
		displayName = strings.TrimPrefix(displayName, "terraform-")
		d := displayMeta{
			Name:        repo.GetName(),
			DisplayName: displayName,
			URL:         repo.GetHTMLURL(),
			Stars:       strconv.Itoa(repo.GetStargazersCount()),
			CreatedAt:   repo.GetCreatedAt().Format(renderTimeformat),
			Description: repo.GetDescription(),
			Labels:      repo.Topics,
		}

		// gh topics to categories
		parsedCategories := []string{}
		for _, topic := range repo.Topics {
			p, exists := topicToCategory[topic]
			if exists {
				parsedCategories = append(parsedCategories, p)
			}
			if topic == e2eLabel {
				d.IsE2E = true
			}
		}
		if len(parsedCategories) > 0 {
			sort.Strings(parsedCategories)
			d.Categories = strings.Join(parsedCategories, ", ")
		}
		dm = append(dm, d)
	}
	return dm
}

// docSort sorts displayMeta surfacing e2e blueprints first for documentation.
func docSort(dm []displayMeta) []displayMeta {
	sort.SliceStable(dm, func(i, j int) bool {
		if dm[i].IsE2E && dm[j].IsE2E {
			return dm[i].DisplayName < dm[j].DisplayName
		}
		return dm[i].IsE2E
	})
	return dm
}

// renderDocHTML renders html for documentation.
func renderDocHTML(dm []displayMeta) string {
	htmlTmpl, err := template.New("htmlDoc").Parse(htmlTemplate)
	if err != nil {
		return fmt.Sprintf("error parsing template: %v", err)
	}
	var tpl bytes.Buffer
	err = htmlTmpl.Execute(&tpl, docSort(dm))
	if err != nil {
		return fmt.Sprintf("error executing template: %v", err)
	}
	return tpl.String()
}
