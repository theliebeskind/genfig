package writers

var (
	indent   = "  " // default is two spaces
	maxLevel = 5    // default is 5 maximum levels of recursion
	nl       = "\n"
)

// SetIndent sets the indent to be used by the writers
// to indent recursive data
func SetIndent(s string) {
	indent = s
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
