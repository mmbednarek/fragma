package util

import (
	"os"
	"testing"
)

func TestTable_Print(t1 *testing.T) {
	table := NewTable()
	table.Add("a", "14")
	table.Add("a", "2")
	table.Add("a", "3")
	table.Add("b", "5")
	table.Add("b", "6435345345")
	table.Add("c", "dd")
	table.Print(os.Stdout)
}
