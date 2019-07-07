package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theliebeskind/go-genfig/types"
	"github.com/theliebeskind/go-genfig/util"
	"github.com/theliebeskind/go-genfig/writers"
)

const (
	fixturesDir = "../fixtures/"
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
			v = append(v.([]string), defaultConfigFilePrefix+e+".go")
		}
		return v
	}, []string{}).([]string)
	goFiles = append(goFiles, defaultSchemaFilename)
	goFiles = append(goFiles, defaultEnvsFilename)

	type args struct {
		files []string
		p     types.Params
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"nonexisting file(s)", args{[]string{"nope.yml"}, types.Params{}}, true},
		{"not a config file", args{[]string{"../notaconfig.txt"}, types.Params{}}, true},
		{"duplicate env files", args{duplicateConfigFiles, types.Params{}}, true},
		{"no default", args{configFilesWithoutDefault, types.Params{}}, true},
		{"too many levels", args{tooManyLevelsConfigFiles, types.Params{DefaultEnv: "default.with.too.many.levels"}}, true},
		{"additional field(s) to default", args{goodConfigFiles, types.Params{DefaultEnv: "local"}}, true},
		{"non conformant value to default", args{nonconformatnConfigFiles, types.Params{}}, true},
		{"existing files, no dir", args{goodConfigFiles, types.Params{}}, false},
		{"existing files with dir", args{goodConfigFiles, types.Params{Dir: filepath.Clean("../out")}}, false},
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

func Test_WriteSchema(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]interface{}
		contains   []string
		wantSchema types.Schema
		wantErr    bool
	}{
		{"empty", map[string]interface{}{}, []string{}, types.Schema{}, false},
		{"simple string", map[string]interface{}{"a": "b"}, []string{"A string"}, types.Schema{}, false},
		{"simple int", map[string]interface{}{"a": 1}, []string{"A int64"}, types.Schema{}, false},
		{"simple bool", map[string]interface{}{"a": true}, []string{"A bool"}, types.Schema{}, false},
		{"int array", map[string]interface{}{"a": []int{1, 2, 3}}, []string{"A []int"}, types.Schema{}, false},
		{"empty interface array", map[string]interface{}{"a": []interface{}{}}, []string{"A []interface {}"}, types.Schema{}, false},
		{"mixed interface array", map[string]interface{}{"a": []interface{}{1, ""}}, []string{"A []interface {}"}, types.Schema{}, false},
		{"int interface array", map[string]interface{}{"a": []interface{}{1, 2}}, []string{"A []int64"}, types.Schema{}, false},
		{"string interface array", map[string]interface{}{"a": []interface{}{"a", "b"}}, []string{"A []string"}, types.Schema{}, false},
		{"map", map[string]interface{}{"a": map[string]interface{}{"b": 1}}, []string{"A struct {", "B int"}, types.Schema{}, false},
		{"iface key map", map[string]interface{}{"a": map[interface{}]interface{}{"b": 1}}, []string{"A struct {", "B int"}, types.Schema{}, false},
		{"map of map", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, []string{"A struct {", "B struct {", "C int"}, types.Schema{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &strings.Builder{}
			_, err := writers.WriteAndReturnSchema(s, tt.config)
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

func Test_WriteConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		def      map[string]interface{}
		contains []string
		wantErr  bool
	}{
		{"empty", map[string]interface{}{}, nil, []string{}, false},
		{"simple string", map[string]interface{}{"a": "b"}, nil, []string{"A: \"b\""}, false},
		{"simple string with default", map[string]interface{}{"a": "b"}, map[string]interface{}{"a": "def"}, []string{"A: \"b\""}, false},
		{"simple int", map[string]interface{}{"a": 1}, nil, []string{"A: 1"}, false},
		{"simple bool", map[string]interface{}{"a": true}, nil, []string{"A: true"}, false},
		{"int array", map[string]interface{}{"a": []int{1, 2, 3}}, nil, []string{"A: []int"}, false},
		{"empy interface array", map[string]interface{}{"a": []interface{}{}}, nil, []string{"A: []interface {}{}"}, false},
		{"mixed interface array", map[string]interface{}{"a": []interface{}{1, ""}}, nil, []string{"A: []interface {}{1, \"\"}"}, false},
		{"int interface array", map[string]interface{}{"a": []interface{}{1, 2}}, nil, []string{"A: []int{1, 2}"}, false},
		{"string interface array", map[string]interface{}{"a": []interface{}{"a", "b"}}, nil, []string{"A: []string"}, false},
		{"map", map[string]interface{}{"a": map[string]interface{}{"b": 1}}, nil, []string{"A: ConfigA{", "B: 1"}, false},
		{"map of map", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, nil, []string{"A: ConfigA{", "B: ConfigAB{", "C: 1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &strings.Builder{}
			def := tt.def
			if def == nil {
				def = tt.config
			}
			err := writers.WriteConfig(s, types.SchemaMap{
				"ConfigA":   types.Schema{},
				"ConfigAB":  types.Schema{},
				"ConfigABC": types.Schema{},
			}, tt.config, def, "test")
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

func Benchmark_WriteConfigValue(b *testing.B) {
	w := util.NoopWriter{}
	s := types.SchemaMap{
		"ConfigA":   types.Schema{},
		"ConfigAB":  types.Schema{},
		"ConfigABC": types.Schema{},
		"ConfigABD": types.Schema{},
		"ConfigABE": types.Schema{},
	}
	m := map[string]interface{}{"a": map[interface{}]interface{}{"b": map[string]interface{}{"c": []interface{}{1}, "d": "s", "e": 1}}}
	e := map[string]interface{}{}
	for n := 0; n < b.N; n++ {
		writers.WriteConfigValue(w, "Config", m, e, s, 0)
	}
}

func Benchmark_WriteSchemaType(b *testing.B) {
	w := util.NoopWriter{}
	s := types.SchemaMap{}
	m := map[string]interface{}{"a": map[interface{}]interface{}{"b0": 1, "b": map[string]interface{}{"c": []interface{}{1}, "d": "s", "e": 1}}}
	for n := 0; n < b.N; n++ {
		writers.WriteSchemaType(w, "Config", m, s, 0)
	}
}
