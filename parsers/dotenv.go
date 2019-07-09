package parsers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	mergo "github.com/imdario/mergo"
	"github.com/theliebeskind/genfig/util"
)

// DotenvStrategy parses yaml and json files
type DotenvStrategy struct {
}

// Parse of DotenvStrategy parses yaml and json files into Parsing result
func (s *DotenvStrategy) Parse(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data")
	}

	r := map[string]interface{}{}

	scanner := bufio.NewScanner(bytes.NewBuffer(data))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		var kv []string
		if kv = strings.SplitN(line, "=", 2); len(kv) != 2 {
			if kv = strings.SplitN(line, ":", 2); len(kv) != 2 {
				return nil, fmt.Errorf("Invalid dotenv line: '%s'", line)
			}
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		keys := strings.Split(strings.ToLower(k), "_")

		var item interface{} = r
		for i, key := range keys {
			if item, found := item.(map[string]interface{})[key]; found {
				switch item.(type) {
				case map[string]interface{}:
					if len(keys) == i+1 {
						return nil, fmt.Errorf("Key '%s' is already present with differnt type (old: map, new: basic)", keys)
					}
				default:
					return nil, fmt.Errorf("Key '%s' is already present with different type (old: basic, new: map)", keys)
				}
			}
		}

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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return r, nil
}
