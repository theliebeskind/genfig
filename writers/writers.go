package writers

import "strings"

var (
	indent   = "\t" // default is two spaces
	maxLevel = 5    // default is 5 maximum levels of recursion
	nl       = "\n" // default is *nix new line

	indents = strings.Repeat(indent, maxLevel+1)
)

// SetIndent sets the indent to be used by the writers
// to indent recursive data
func SetIndent(s string) {
	indent = s
	indents = strings.Repeat(indent, maxLevel+1)
}

// SetMaxLevel sets the maximum level of recursion;
// If one configuration exceeds this maximum level,
// the generation fails
func SetMaxLevel(l int) {
	maxLevel = l
}

// SetNewline sets the new line to be used by the writers
func SetNewline(s string) {
	nl = s
}
