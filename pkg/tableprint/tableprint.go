package tableprint

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

func Print0[T any](header []string, records []T, rowFormat func(r any) []any, opts ...tablewriter.Option) {
	opts = append(opts,
		tablewriter.WithHeader(header),
		tablewriter.WithTrimSpace(tw.Off),
	)

	_table := tablewriter.NewTable(
		os.Stdout,
		opts...,
	)

	for _, r := range records {
		_ = _table.Append(rowFormat(r)...)
	}
	_ = _table.Render()
}

func Print[T any](header []any, records []T, rowFormat func(r any) []any) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	for _, r := range records {
		t.AppendRow(rowFormat(r))
	}
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	//t.SetStyle(table.Style{
	//	Name:   "CustomStyle",
	//	Box:    table.StyleBoxDefault,
	//	Format: table.FormatOptionsDefault,
	//})
	t.Render()
}

func JustSelf(r any) []any {
	r0 := r.([]any)
	return r0
}
