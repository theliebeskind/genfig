package parsers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	mergo "github.com/imdario/mergo"
	"github.com/thclodes/genfig/util"
)

var (
	allowedKVSeparators  = []string{"=", ":"}
	allowedEnvSeparators = []string{"_", "-"}
	keyRegex             = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)
)

// DotenvStrategy parses .env files.
// Emply lines are ignored.
// Allowed key-value separators are only "=" and ":".
// Supported values are: string, int64, float64, bool and json arrays,
// e.g. `["a", "b", "c"]` ([]string) or `["a", 1, true]` ([]interface {}).
// Key can be nests by eithe one of the allowed env separators, e.g. `DB_NAME` or `SERVER-HOST`
// Comments are only allowed in seperate lines.
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
		line := strings.TrimSpace(scanner.Text())
		// ignore empty or comment lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// try to split line into key and value by allowed separators
		var kv []string
		for i, sep := range allowedKVSeparators {
			if kv = strings.SplitN(line, sep, 2); len(kv) == 2 {
				break
			}
			if i == len(allowedKVSeparators)-1 {
				return nil, fmt.Errorf("Invalid dotenv line: '%s'", line)
			}
		}

		// trim key and value
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])

		var keys []string
		for _, sep := range allowedEnvSeparators {
			if keys = strings.Split(strings.ToLower(k), sep); len(keys) > 1 {
				break
			}
		}

		var item interface{} = r
		for i, key := range keys {
			if !keyRegex.MatchString(key) {
				return nil, fmt.Errorf("Key '%s' is not valid", key)
			}
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
