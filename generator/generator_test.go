package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thlcodes/genfig/models"
	"github.com/thlcodes/genfig/util"
)

var (
	fixturesDir, _ = filepath.Abs("../fixtures/")
)

func Test_Generate(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "genfig")
	defer os.RemoveAll(tmpDir)

	cwd, _ := os.Getwd()
	os.MkdirAll(filepath.Join(tmpDir, "workdir"), 0777)
	os.Chdir(filepath.Join(tmpDir, "workdir"))
	defer os.Chdir(cwd)

	configsDir := filepath.Join(fixturesDir, "configs/")

	configFilesWithoutDefault := []string{
		configsDir + "/.env.local",
		configsDir + "/development.local.toml",
		configsDir + "/development.yaml",
		configsDir + "/production.json",
		configsDir + "/test.json",
	}
	goodConfigFiles := append(configFilesWithoutDefault, configsDir+"/default.yml")
	duplicateConfigFiles := []string{configsDir + "/local.yml", configsDir + "/local.yml"}
	nonconformatnConfigFiles := []string{configsDir + "/default.yml", configsDir + "/nonconformant.yml"}
	tooManyLevelsConfigFiles := []string{configsDir + "/default.with.too.many.levels.yml"}

	goFiles := util.ReduceStrings(goodConfigFiles, func(r interface{}, s string) interface{} {
		e, _ := parseFilename(s)
		v := r
		if e != "" && e != "default" {
			if e == "test" {
				e = "test_"
			}
			v = append(v.([]string), defaultConfigFilePrefix+e+".go")
		}
		return v
	}, []string{}).([]string)
	goFiles = append(goFiles, defaultSchemaFilename)
	goFiles = append(goFiles, defaultEnvsFilename)

	type args struct {
		files []string
		p     models.Params
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"nonexisting file(s)", args{[]string{"nope.yml"}, models.Params{}}, true},
		{"not a config file", args{[]string{configsDir + "/notaconfig.txt"}, models.Params{}}, true},
		{"duplicate env files", args{duplicateConfigFiles, models.Params{}}, true},
		{"no default", args{configFilesWithoutDefault, models.Params{}}, true},
		{"too many levels", args{tooManyLevelsConfigFiles, models.Params{DefaultEnv: "default.with.too.many.levels"}}, true},
		{"additional field(s) to default", args{goodConfigFiles, models.Params{DefaultEnv: "local"}}, true},
		{"non conformant value to default", args{nonconformatnConfigFiles, models.Params{}}, true},
		{"existing files, no dir", args{goodConfigFiles, models.Params{}}, false},
		{"existing files with dir", args{goodConfigFiles, models.Params{Dir: "config"}}, false},
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
