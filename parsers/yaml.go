package strategies

import (
	"errors"

	yaml "gopkg.in/yaml.v2"
)

// YamlStrategy parses yaml and json files
type YamlStrategy struct {
}

// Parse of YamlStrategy parses yaml and json files into Parsing result
func (s *YamlStrategy) Parse(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data")
	}
	r := map[string]interface{}{}

	err := yaml.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
