package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theliebeskind/genfig/types"
	"github.com/theliebeskind/genfig/util"
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
