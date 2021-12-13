package bptest

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func newTable() table.Writer {
	tw := table.NewWriter()
	tw.Style().Color.Header = text.Colors{text.FgGreen}
	tw.SetColumnConfigs(
		[]table.ColumnConfig{
			{Number: 1, Colors: text.Colors{text.FgYellow}},
		},
	)
	tw.Style().Options.DrawBorder = false
	tw.SetOutputMirror(os.Stdout)
	return tw
}
