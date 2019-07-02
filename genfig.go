package genfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/theliebeskind/genfig/strategies"
	"github.com/theliebeskind/genfig/util"
)

const (
	defaultEnvName          = "default"
	defaultSchemaFilename   = "schema.go"
	defaultEnvsFilename     = "envs.go"
	defaulGenfigFilename    = "genfig.go"
	defaultConfigFilePrefix = "env_"
	defaultPackage          = "config"
	defaultCmd              = "genfig"
	defaultIndent           = "  "
	defaultNewline          = "\n"
)

const (
	maxLevel = 5
)

var (
	ymlStrategy    = strategies.YamlStrategy{}
	tomlStrategy   = strategies.TomlStrategy{}
	dotenvStrategy = strategies.DotenvStrategy{}

	nl = defaultNewline
)

var (
	allowedExtensions = []string{"\\.yml", "\\.yaml", "\\.json", "\\.toml"}
	allowedPrefixes   = []string{"\\.env"}
	strategiesMap     = map[string]strategies.ParsingStrategy{
		"yml":    &ymlStrategy,
		"json":   &ymlStrategy,
		"toml":   &tomlStrategy,
		"dotenv": &dotenvStrategy,
	}
	envReStr = `((?:` + strings.Join(allowedPrefixes, "|") + `)\.([\w\.]+))|(([\w\.]+)(` + strings.Join(allowedExtensions, "|") + `))`
	envRe    = regexp.MustCompile(envReStr)
)

// Schema defines the schema
type Schema struct {
	IsStruct bool
	Content  string
}

// SchemaMap aliases as string-map of bytes
type SchemaMap map[string]Schema

// Params for the Generate func as struct,
// empty values are default values, so can be passed empty
type Params struct {
	Dir        string
	DefaultEnv string
	MergeFiles bool
}

// Generate generates the go config files
func Generate(files []string, params Params) ([]string, error) {
	var err error
	if len(files) == 0 {
		return nil, errors.New("No files to generate from")
	}

	envs := []string{}
	envMap := make(map[string]map[string]interface{})
	fileMap := make(map[string]string)

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return nil, err
		}

		env, typ := parseFilename(filepath.Base(f))
		if env == "" {
			continue
		}
		if _, exists := strategiesMap[typ]; !exists {
			continue
		}
		if _, exists := envMap[env]; exists {
			return nil, fmt.Errorf("Environment '%s' does already exist", env)
		}
		var err error
		envMap[env], err = parseFile(f, strategiesMap[typ])
		if err != nil {
			return nil, err
		}
		fileMap[env] = f
	}

	if len(envMap) == 0 {
		return nil, errors.New("No suitable config files found")
	}

	if params.DefaultEnv == "" {
		params.DefaultEnv = defaultEnvName
	}

	var defaultEnv map[string]interface{}
	var hasDefault bool
	if defaultEnv, hasDefault = envMap[params.DefaultEnv]; !hasDefault {
		return nil, errors.New("Missing default config")
	}

	if err := os.MkdirAll(params.Dir, 0777); params.Dir != "" && err != nil {
		return nil, err
	}

	var schema SchemaMap
	schemaFileName := filepath.Join(params.Dir, defaultSchemaFilename)
	source := fmt.Sprintf("%s (schema built from '%s')", defaultCmd, filepath.Base(fileMap[params.DefaultEnv]))
	if err := func() (err error) {
		var f *os.File
		defer func() {
			if f != nil {
				_ = f.Close()
			}
		}()
		if f, err = os.Create(schemaFileName); err != nil {
			return err
		} else if err = writeHeader(f, defaultPackage, source); err != nil {
			return err
		} else if schema, err = writeAndReturnSchema(f, defaultEnv); err != nil {
			return err
		}
		return
	}(); err != nil {
		return nil, err
	}

	gofiles := []string{schemaFileName}

	for env, data := range envMap {
		out := defaultConfigFilePrefix + env + ".go"
		path := filepath.Join(params.Dir, out)
		source := fmt.Sprintf("%s (config built by merging '%s' and '%s')", defaultCmd, filepath.Base(fileMap[params.DefaultEnv]), filepath.Base(fileMap[env]))
		name := strings.ReplaceAll(strings.Title(strings.ReplaceAll(env, "_", ".")), ".", "")
		envs = append(envs, name)

		// Check of schema of this config does conform the the global schema
		// If is has additional fields or fields with different schema themselves,
		// it fails
		var configSchema SchemaMap
		if configSchema, err = writeAndReturnSchema(util.NoopWriter{}, data); err != nil {
			return nil, err
		}
		for k, s := range configSchema {
			if s.IsStruct {
				continue
			}
			if _, exists := schema[k]; !exists || s.Content != schema[k].Content {
				return nil, fmt.Errorf("%s has at leas one non-conformant field: '%s': %s != '%s'", fileMap[env], k, s.Content, schema[k].Content)
			}
		}

		if err := func() (err error) {
			var f *os.File
			defer func() {
				if f != nil {
					_ = f.Close()
				}
			}()
			if f, err = os.Create(path); err != nil {
				return err
			} else if err = writeHeader(f, defaultPackage, source); err != nil {
				return err
			} else if err = writeConfig(f, schema, data, defaultEnv, name); err != nil {
				return err
			}
			return
		}(); err != nil {
			return nil, err
		}

		gofiles = append(gofiles, path)
	}

	envsFileName := filepath.Join(params.Dir, defaultEnvsFilename)
	if err := func() (err error) {
		var f *os.File
		defer func() {
			if f != nil {
				_ = f.Close()
			}
		}()
		if f, err = os.Create(envsFileName); err != nil {
			return err
		} else if err = writeHeader(f, defaultPackage, defaultCmd); err != nil {
			return err
		} else if err = writeEnvs(f, envs); err != nil {
			return err
		}
		return
	}(); err != nil {
		return nil, err
	}
	gofiles = append(gofiles, envsFileName)

	return gofiles, nil
}

func parseFile(f string, s strategies.ParsingStrategy) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return s.Parse(data)
}

func parseFilename(f string) (string, string) {
	typ := filepath.Ext(f)
	if len(typ) == 0 {
		return "", ""
	}
	typ = typ[1:]
	if typ == "yaml" {
		typ = "yml"
	} else if strings.HasPrefix(f, ".env") {
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
