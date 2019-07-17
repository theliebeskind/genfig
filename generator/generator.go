package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/theliebeskind/genfig/writers"

	"github.com/theliebeskind/genfig/types"

	"github.com/theliebeskind/genfig/parsers"
	"github.com/theliebeskind/genfig/util"
)

const (
	defaultEnvName          = "default"
	defaultSchemaFilename   = "schema.go"
	defaultEnvsFilename     = "envs.go"
	defaultInitFilename     = "init.go"
	defaulGenfigFilename    = "genfig.go"
	defaultConfigFilePrefix = "env_"
	defaultPackage          = "config"
	defaultCmd              = "genfig"
)

var (
	ymlStrategy    = parsers.YamlStrategy{}
	tomlStrategy   = parsers.TomlStrategy{}
	dotenvStrategy = parsers.DotenvStrategy{}
)

var (
	allowedExtensions = []string{"\\.yml", "\\.yaml", "\\.json", "\\.toml"}
	allowedPrefixes   = []string{"\\.env"}
	parsersMap        = map[string]parsers.ParsingStrategy{
		"yml":    &ymlStrategy,
		"json":   &ymlStrategy,
		"toml":   &tomlStrategy,
		"dotenv": &dotenvStrategy,
	}
	envReStr = `((?:` + strings.Join(allowedPrefixes, "|") + `)\.([\w\.]+))|(([\w\.]+)(` + strings.Join(allowedExtensions, "|") + `))`
	envRe    = regexp.MustCompile(envReStr)
)

// Generate generates the go config files
func Generate(files []string, params types.Params) ([]string, error) {
	var err error
	if len(files) == 0 {
		return nil, errors.New("No files to generate from")
	}

	if !filepath.IsAbs(params.Dir) {
		wd, _ := os.Getwd()
		params.Dir = filepath.Join(wd, params.Dir)
	}

	envs := map[string]string{}
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
		if _, exists := parsersMap[typ]; !exists {
			continue
		}
		if _, exists := envMap[env]; exists {
			return nil, fmt.Errorf("Environment '%s' does already exist", env)
		}
		var err error
		envMap[env], err = parseFile(f, parsersMap[typ])
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

	gofiles := []string{}

	// write schemafile
	var schema types.SchemaMap
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
		} else if err = writers.WriteHeader(f, defaultPackage, source); err != nil {
			return err
		} else if schema, err = writers.WriteAndReturnSchema(f, defaultEnv); err != nil {
			return err
		}
		return
	}(); err != nil {
		return nil, err
	}
	gofiles = append(gofiles, schemaFileName)

	// write config files
	for env, data := range envMap {
		out := defaultConfigFilePrefix + env + ".go"
		path := filepath.Join(params.Dir, out)
		source := fmt.Sprintf("%s (config built by merging '%s' and '%s')", defaultCmd, filepath.Base(fileMap[params.DefaultEnv]), filepath.Base(fileMap[env]))
		name := strings.ReplaceAll(strings.Title(strings.ReplaceAll(env, "_", ".")), ".", "")
		envs[env] = name

		// Check of schema of this config does conform the the global schema
		// If is has additional fields or fields with different schema themselves,
		// it fails
		var configSchema types.SchemaMap
		if configSchema, err = writers.WriteAndReturnSchema(util.NoopWriter{}, data); err != nil {
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
			} else if err = writers.WriteHeader(f, defaultPackage, source); err != nil {
				return err
			} else if err = writers.WriteConfig(f, schema, data, defaultEnv, name); err != nil {
				return err
			}
			return
		}(); err != nil {
			return nil, err
		}

		gofiles = append(gofiles, path)
	}

	// write env file
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
		} else if err = writers.WriteHeader(f, defaultPackage, defaultCmd); err != nil {
			return err
		} else if err = writers.WriteEnvs(f, envs); err != nil {
			return err
		}
		return
	}(); err != nil {
		return nil, err
	}
	gofiles = append(gofiles, envsFileName)

	pluginCalls := map[string]string{}
	// write plugins files
	var pfiles []string
	if pfiles, err = writers.WritePlugins(schema, params.Dir, defaultPackage, defaultCmd, pluginCalls); err != nil {
		return nil, err
	}
	gofiles = append(gofiles, pfiles...)

	// write init file
	initFileName := filepath.Join(params.Dir, defaultInitFilename)
	if err := func() (err error) {
		var f *os.File
		defer func() {
			if f != nil {
				_ = f.Close()
			}
		}()
		if f, err = os.Create(initFileName); err != nil {
			return err
		} else if err = writers.WriteHeader(f, defaultPackage, defaultCmd); err != nil {
			return err
		} else if err = writers.WriteInit(f, pluginCalls); err != nil {
			return err
		}
		return
	}(); err != nil {
		return nil, err
	}
	gofiles = append(gofiles, schemaFileName)

	return gofiles, nil
}

func parseFile(f string, s parsers.ParsingStrategy) (map[string]interface{}, error) {
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
