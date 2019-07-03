package strategies

import (
	"errors"

	toml "github.com/BurntSushi/toml"
)

// TomlStrategy parses yaml and json files
type TomlStrategy struct {
}

// Parse of TomlStrategy parses yaml and json files into Parsing result
func (s *TomlStrategy) Parse(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data")
	}
	r := map[string]interface{}{}

	err := toml.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
