package genfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theliebeskind/genfig/util"
)

const (
	fixturesDir = "./fixtures/"
)

func Test_Generate(t *testing.T) {
	cwd, _ := os.Getwd()
	os.Chdir(filepath.Join(fixturesDir, "cwd"))
	defer os.Chdir(cwd)

	configFilesWithoutDefault := []string{
		"../.env.local",
		"../development.local.toml",
		"../development.yaml",
		"../production.json",
	}
	goodConfigFiles := append(configFilesWithoutDefault, "../default.yml")
	duplicateConfigFiles := []string{"../local.yml", "../local.yml"}
	nonconformatnConfigFiles := []string{"../default.yml", "../nonconformant.yml"}
	tooManyLevelsConfigFiles := []string{"../default.with.too.many.levels.yml"}

	goFiles := util.ReduceStrings(goodConfigFiles, func(r interface{}, s string) interface{} {
		e, _ := parseFilename(s)
		v := r
		if e != "" && e != "default" {
			v = append(v.([]string), e+".go")
		}
		return v
	}, []string{}).([]string)
	goFiles = append(goFiles, "genfig_schema.go")

	type args struct {
		files []string
		p     Params
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"nonexisting file(s)", args{[]string{"nope.yml"}, Params{}}, true},
		{"not a config file", args{[]string{"../notaconfig.txt"}, Params{}}, true},
		{"duplicate env files", args{duplicateConfigFiles, Params{}}, true},
		{"no default", args{configFilesWithoutDefault, Params{}}, true},
		{"too many levels", args{tooManyLevelsConfigFiles, Params{DefaultEnv: "default.with.too.many.levels"}}, true},
		{"additional field(s) to default", args{goodConfigFiles, Params{DefaultEnv: "local"}}, true},
		{"non conformant value to default", args{nonconformatnConfigFiles, Params{}}, true},
		{"existing files, no dir", args{goodConfigFiles, Params{}}, false},
		{"existing files with dir", args{goodConfigFiles, Params{Dir: filepath.Clean("../out")}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.args.p.Dir
			if dir != "" && !strings.HasSuffix(dir, "/") {
				dir += "/"
			}
			if dir != "" {
				os.RemoveAll(dir)
			} else {
				for _, gofile := range goFiles {
					_ = os.Remove(dir + gofile)
				}
			}

			_, err := Generate(tt.args.files, tt.args.p)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			for _, gofile := range goFiles {
				assert.FileExists(t, filepath.Clean(dir+gofile))
			}
		})
	}
}

func Test_parseFilename(t *testing.T) {
	env := "environ.local"
	tests := []struct {
		name    string
		f       string
		wantEnv string
		wantTyp string
	}{
		{"empty", "", "", ""},
		{"yml", env + ".yml", env, "yml"},
		{"yaml", env + ".yaml", env, "yml"},
		{"json", env + ".json", env, "json"},
		{"toml", env + ".toml", env, "toml"},
		{"dotenv", ".env." + env, env, "dotenv"},
		{"noext", env, "", "local"},
		{"invalidext", env + ".bla", "", "bla"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, typ := parseFilename(tt.f)
			assert.Equal(t, tt.wantEnv, env)
			assert.Equal(t, tt.wantTyp, typ)
		})
	}
}

func Test_writeSchema(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]interface{}
		contains   []string
		wantSchema Schema
		wantErr    bool
	}{
		{"empty", map[string]interface{}{}, []string{}, Schema{}, false},
		{"simple string", map[string]interface{}{"a": "b"}, []string{"A string"}, Schema{}, false},
		{"simple int", map[string]interface{}{"a": 1}, []string{"A int64"}, Schema{}, false},
		{"simple bool", map[string]interface{}{"a": true}, []string{"A bool"}, Schema{}, false},
		{"int array", map[string]interface{}{"a": []int{1, 2, 3}}, []string{"A []int"}, Schema{}, false},
		{"empy interface array", map[string]interface{}{"a": []interface{}{}}, []string{"A []interface {}"}, Schema{}, false},
		{"mixed interface array", map[string]interface{}{"a": []interface{}{1, ""}}, []string{"A []interface {}"}, Schema{}, false},
		{"int interface array", map[string]interface{}{"a": []interface{}{1, 2}}, []string{"A []int"}, Schema{}, false},
		{"string interface array", map[string]interface{}{"a": []interface{}{"a", "b"}}, []string{"A []string"}, Schema{}, false},
		{"map", map[string]interface{}{"a": map[string]interface{}{"b": 1}}, []string{"A struct {", "B int"}, Schema{}, false},
		{"iface key map", map[string]interface{}{"a": map[interface{}]interface{}{"b": 1}}, []string{"A struct {", "B int"}, Schema{}, false},
		{"map of map", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, []string{"A struct {", "B struct {", "C int"}, Schema{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &strings.Builder{}
			_, err := writeAndReturnSchema(s, tt.config)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			got := s.String()
			for _, c := range tt.contains {
				assert.Contains(t, got, c)
			}
		})
	}
}

func Test_writeConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		contains []string
		wantErr  bool
	}{
		{"empty", map[string]interface{}{}, []string{}, false},
		{"simple string", map[string]interface{}{"a": "b"}, []string{"A: \"b\""}, false},
		{"simple int", map[string]interface{}{"a": 1}, []string{"A: 1"}, false},
		{"simple bool", map[string]interface{}{"a": true}, []string{"A: true"}, false},
		{"int array", map[string]interface{}{"a": []int{1, 2, 3}}, []string{"A: []int"}, false},
		{"empy interface array", map[string]interface{}{"a": []interface{}{}}, []string{"A: []interface {}{}"}, false},
		{"mixed interface array", map[string]interface{}{"a": []interface{}{1, ""}}, []string{"A: []interface {}{1, \"\"}"}, false},
		{"int interface array", map[string]interface{}{"a": []interface{}{1, 2}}, []string{"A: []int{1, 2}"}, false},
		{"string interface array", map[string]interface{}{"a": []interface{}{"a", "b"}}, []string{"A: []string"}, false},
		{"map", map[string]interface{}{"a": map[string]interface{}{"b": 1}}, []string{"A: ConfigA{", "B: 1"}, false},
		{"map of map", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, []string{"A: ConfigA{", "B: ConfigAB{", "C: 1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &strings.Builder{}
			err := writeConfig(s, SchemaMap{
				"ConfigA":   Schema{},
				"ConfigAB":  Schema{},
				"ConfigABC": Schema{},
			}, tt.config, "test")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			got := s.String()
			require.NoError(t, err)
			for _, c := range tt.contains {
				assert.Contains(t, got, c)
			}
		})
	}
}
