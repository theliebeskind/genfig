package writers_test

import "github.com/thlcodes/genfig/writers"

const (
	maxLevel = 10
	indent   = "  "
	newLine  = "n"
)

func init() {
	writers.SetIndent(indent)
	writers.SetMaxLevel(maxLevel)
	writers.SetNewline(newLine)
}
