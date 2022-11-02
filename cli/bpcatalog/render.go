package bpcatalog

import (
	"fmt"
	"io"

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
	renderTable      renderFormat = "table"
	renderCSV        renderFormat = "csv"
	renderTimeformat string       = "2006-01-02"
)

var (
	renderFormats = []renderFormat{renderTable, renderCSV}
)

// render writes given repo information in the specified renderFormat to w.
func render(r repos, w io.Writer, format renderFormat) error {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(w)
	tbl.AppendHeader(table.Row{"Repo", "Stars", "Created"})
	for _, repo := range r {
		tbl.AppendRow(table.Row{repo.GetName(), repo.GetStargazersCount(), repo.GetCreatedAt().Format(renderTimeformat)})
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
