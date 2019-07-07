package parsers

import (
	"errors"
	"strings"

	mergo "github.com/imdario/mergo"
	dotenv "github.com/joho/godotenv"
	"github.com/theliebeskind/go-genfig/util"
)

// DotenvStrategy parses yaml and json files
type DotenvStrategy struct {
}

// Parse of DotenvStrategy parses yaml and json files into Parsing result
func (s *DotenvStrategy) Parse(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data")
	}

	envs, err := dotenv.Unmarshal(string(data))
	if err != nil {
		return nil, err
	}

	r := map[string]interface{}{}
	for k, v := range envs {
		keys := strings.Split(strings.ToLower(k), "_")
		util.ReverseStrings(keys)
		tmp := map[string]interface{}{}

		for i, k := range keys {
			if i == 0 {
				tmp[k] = util.ParseString(v)
				continue
			}
			tmp = map[string]interface{}{k: tmp}
		}

		mergo.Map(&r, tmp)
	}

	return r, nil
}
