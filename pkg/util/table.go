package util

import (
	"io"
	"strings"
)

type Column struct {
	Label string
	Data  []string
}

func (c *Column) Width() int {
	max := len(c.Label)
	for _, value := range c.Data {
		if len(value) > max {
			max = len(value)
		}
	}
	return max
}

func (c *Column) Height() int {
	return len(c.Data)
}

type Table struct {
	Columns    map[string]*Column
	ColumOrder []string
}

func NewTable() *Table {
	return &Table{Columns: map[string]*Column{}, ColumOrder: []string{}}
}

func (t *Table) Add(column string, value string) {
	col, ok := t.Columns[column]
	if !ok {
		t.ColumOrder = append(t.ColumOrder, column)
		t.Columns[column] = &Column{
			Label: strings.ToUpper(column),
			Data:  []string{value},
		}
		return
	}
	col.Data = append(col.Data, value)
}

func (t *Table) Print(out io.Writer) {
	height := 0
	widths := map[string]int{}
	for _, name := range t.ColumOrder {
		col := t.Columns[name]
		widths[name] = col.Width()
		if height < col.Height() {
			height = col.Height()
		}
	}

	for _, name := range t.ColumOrder {
		col := t.Columns[name]
		width := widths[name]
		out.Write([]byte("\033[1;33m"))
		out.Write([]byte(col.Label))
		for j := 0; j < (2 + width - len(col.Label)); j++ {
			out.Write([]byte{' '})
		}
	}
	out.Write([]byte{'\n'})

	out.Write([]byte("\033[0m"))

	for i := 0; i < height; i++ {
		for _, name := range t.ColumOrder {
			col := t.Columns[name]
			width := widths[name]
			if i >= col.Height() {
				for j := 0; j < (2 + width); j++ {
					out.Write([]byte{' '})
				}
				continue
			}

			out.Write([]byte(col.Data[i]))
			for j := 0; j < (2 + width - len(col.Data[i])); j++ {
				out.Write([]byte{' '})
			}
		}
		out.Write([]byte{'\n'})
	}
}
