package strategies

import (
	"errors"

	toml "github.com/BurntSushi/toml"
)

// TomlStrategy parses yaml and json files
type TomlStrategy struct {
}

// Parse of TomlStrategy parses yaml and json files into Parsing result
func (s *TomlStrategy) Parse(data []byte) (ParsingResult, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data")
	}
	r := ParsingResult{}

	_, err := toml.Decode(string(data), &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
