package strategies

// ParsingStrategy interface
type ParsingStrategy interface {
	Parse(data []byte) (map[string]interface{}, error)
}
