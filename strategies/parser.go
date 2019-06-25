package strategies

// ParsingResult is an alias for generic dict
type ParsingResult map[string]interface{}

// ParsingStrategy interface
type ParsingStrategy interface {
	Parse(data []byte) (ParsingResult, error)
}
