package bptest

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func newTable(header ...interface{}) table.Table {
	headerFmt := color.New(color.FgGreen).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New(header...)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
