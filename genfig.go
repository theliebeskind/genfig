package genfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	allowedExtensions = []string{"yml", "yaml", "toml", "json"}
	envRe             = regexp.MustCompile(`(\.env\.([\w\.]+))|(([\w\.]+)\.(` + strings.Join(allowedExtensions, "|") + `))`)
)

// Data is an alias
type Data map[string]interface{}

// Generate generates the go config files
func Generate(files []string, dir string) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("No files to generate from")
	}

	envs := make(map[string]Data)

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return nil, err
		}

		env := extractEnviron(f)
		if _, exists := envs[env]; exists {
			return nil, fmt.Errorf("Environment '%s' already exists", env)
		}
		var err error
		envs[env], err = parseFile(f)
		if err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}

	gofiles := make([]string, len(envs))
	i := 0
	for env, data := range envs {
		out := env + ".go"
		path := filepath.Join(dir, out)
		if err := writeData(data, path); err != nil {
			return nil, err
		}
		gofiles[i] = path
		i++
	}
	return gofiles, nil
}

func parseFile(f string) (Data, error) {
	return Data{}, nil
}

func writeData(d Data, to string) error {
	return ioutil.WriteFile(to, []byte{}, 0777)
}

func extractEnviron(f string) string {
	match := envRe.FindAllStringSubmatch(f, 1)
	if len(match) == 0 {
		return ""
	}
	if match[0][2] != "" {
		return match[0][2]
	}
	if match[0][4] != "" {
		return match[0][4]
	}
	return ""
}

func isConfigFile(f string) bool {
	return extractEnviron(f) != ""
}
