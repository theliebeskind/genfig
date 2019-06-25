package strategies

import (
	"errors"

	yaml "gopkg.in/yaml.v2"
)

// YamlStrategy parses yaml and json files
type YamlStrategy struct {
}

// Parse of YamlStrategy parses yaml and json files into Parsing result
func (s *YamlStrategy) Parse(data []byte) (ParsingResult, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data")
	}
	r := ParsingResult{}

	err := yaml.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
