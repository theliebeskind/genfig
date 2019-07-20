package writers_test

import "github.com/theliebeskind/genfig/writers"

const (
	maxLevel = 5
	indent   = "  "
	newLine  = "n"
)

func init() {
	writers.SetIndent(indent)
	writers.SetMaxLevel(maxLevel)
	writers.SetNewline(newLine)
}
