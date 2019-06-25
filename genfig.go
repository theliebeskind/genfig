// Package genfig proveds the genfig methods
//go:generate qtc
package genfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/theliebeskind/genfig/util"

	"github.com/theliebeskind/genfig/strategies"
	"github.com/theliebeskind/genfig/templates"
)

const (
	defaultConfName       = "default"
	defaultSchemaFilename = "config.schema.go"
)

var (
	ymlStrategy    = strategies.YamlStrategy{}
	tomlStrategy   = strategies.TomlStrategy{}
	dotenvStrategy = strategies.DotenvStrategy{}
)

var (
	allowedExtensions = []string{"yml", "yaml", "json", "toml"}
	strategiesMap     = map[string]strategies.ParsingStrategy{
		"yml":    &ymlStrategy,
		"json":   &ymlStrategy,
		"toml":   &tomlStrategy,
		"dotenv": &dotenvStrategy,
	}
	envRe = regexp.MustCompile(`(\.env\.([\w\.]+))|(([\w\.]+)\.(` + strings.Join(allowedExtensions, "|") + `))`)
)

// Generate generates the go config files
func Generate(files []string, dir string) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("No files to generate from")
	}

	envs := make(map[string]strategies.ParsingResult)

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return nil, err
		}

		env, typ := parseFilename(f)
		if env == "" {
			continue
		}
		if _, exists := strategiesMap[typ]; !exists {
			continue
		}
		if _, exists := envs[env]; exists {
			return nil, fmt.Errorf("Environment '%s' already exists", env)
		}
		var err error
		envs[env], err = parseFile(f, strategiesMap[typ])
		if err != nil {
			return nil, err
		}
	}

	if len(envs) == 0 {
		return nil, errors.New("No suitable config files found")
	}

	if _, hasDefault := envs[defaultConfName]; !hasDefault {
		return nil, errors.New("Missing default config")
	}

	if err := os.MkdirAll(dir, 0777); dir != "" && err != nil {
		return nil, err
	}

	if err := writeSchema(envs[defaultConfName], filepath.Join(dir, defaultSchemaFilename)); err != nil {
		return nil, err
	}

	gofiles := make([]string, len(envs))
	i := 0
	for env, data := range envs {
		out := env + ".go"
		path := filepath.Join(dir, out)
		if err := writeConfig(data, path); err != nil {
			return nil, err
		}
		gofiles[i] = path
		i++
	}
	return gofiles, nil
}

func parseFile(f string, s strategies.ParsingStrategy) (strategies.ParsingResult, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return s.Parse(data)
}

func createDefaultSchema(config strategies.ParsingResult) (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = util.RecoverError(r)
			return
		}
	}()
	s = templates.Schema(config)
	return
}

func writeSchema(c strategies.ParsingResult, to string) error {
	s, err := createDefaultSchema(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, []byte(s), 0777)
}

func writeConfig(d strategies.ParsingResult, to string) error {
	data := []byte{}
	return ioutil.WriteFile(to, data, 0777)
}

func parseFilename(f string) (string, string) {
	typ := filepath.Ext(f)
	if len(typ) == 0 {
		return "", ""
	}
	typ = typ[1:]
	if typ == "yaml" {
		typ = "yml"
	} else if strings.HasPrefix(filepath.Base(f), ".env") {
		typ = "dotenv"
	}

	match := envRe.FindAllStringSubmatch(f, 1)
	if len(match) == 0 {
		return "", typ
	}
	if match[0][2] != "" {
		return match[0][2], typ
	}
	if match[0][4] != "" {
		return match[0][4], typ
	}
	return "", typ
}
