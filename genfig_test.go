package genfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theliebeskind/genfig/strategies"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theliebeskind/genfig/util"
)

const (
	fixturesDir = "./fixtures/"
)

func TestGenerate(t *testing.T) {
	cwd, _ := os.Getwd()
	os.Chdir(filepath.Join(fixturesDir, "cwd"))
	defer os.Chdir(cwd)

	configFiles := util.ResolveGlobs("../*")
	configFilesWithoutDefault := util.ResolveGlobs("../*.toml")
	goFiles := util.ReduceStrings(configFiles, func(r interface{}, s string) interface{} {
		e, _ := parseFilename(s)
		v := r
		if e != "" {
			v = append(v.([]string), e+".go")
		}
		return v
	}, []string{}).([]string)
	goFiles = append(goFiles, "config.schema.go")+
	
	duplicateConfigFiles := []string{"../local.yml", "../local.yml"}

	type args struct {
		files []string
		dir   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"nonexisting files", args{[]string{"nope.yml"}, ""}, true},
		{"not a config file", args{[]string{"../notaconfig.txt"}, ""}, true},
		{"duplicate env files", args{duplicateConfigFiles, ""}, true},
		{"no default", args{configFilesWithoutDefault, ""}, true},
		{"existing files, no dir", args{configFiles, ""}, false},
		{"existing files with dir", args{configFiles, filepath.Clean("../out")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate(tt.args.files, tt.args.dir)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.args.dir != "" && !strings.HasSuffix(tt.args.dir, "/") {
				tt.args.dir += "/"
			}
			for _, gofile := range goFiles {
				assert.FileExists(t, filepath.Clean(tt.args.dir+gofile))
				os.Remove(tt.args.dir + gofile)
			}
			if tt.args.dir != "" {
				os.RemoveAll(tt.args.dir)
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

func Test_createDefaultSchema(t *testing.T) {
	tests := []struct {
		name     string
		config   strategies.ParsingResult
		contains []string
		wantErr  bool
	}{
		{"empty", strategies.ParsingResult{}, []string{}, false},
		{"simple string", strategies.ParsingResult{"a": "b"}, []string{"A string"}, false},
		{"simple int", strategies.ParsingResult{"a": int64(1)}, []string{"A int64"}, false},
		{"simple bool", strategies.ParsingResult{"a": true}, []string{"A bool"}, false},
		{"single interface array", strategies.ParsingResult{"a": []interface{}{}}, []string{"A []interface {}"}, false},
		{"single int array", strategies.ParsingResult{"a": []int{1, 2, 3}}, []string{"A []int"}, false},
		{"map", strategies.ParsingResult{"a": map[string]interface{}{"b": 1}}, []string{"A struct {", "B int"}, false},
		{"map of map", strategies.ParsingResult{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, []string{"A struct {", "B struct {", "C int"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createDefaultSchema(tt.config)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			for _, c := range tt.contains {
				assert.Contains(t, got, c)
			}
		})
	}
}
