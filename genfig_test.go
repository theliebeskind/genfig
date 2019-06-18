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

func TestGenerate(t *testing.T) {
	cwd, _ := os.Getwd()
	os.Chdir(filepath.Join(fixturesDir, "cwd"))
	defer os.Chdir(cwd)

	configFiles := util.ResolveGlobs("../*")
	goFiles := util.MapString(configFiles, func(s string) (f string) {
		f = extractEnviron(s) + ".go"
		return
	})
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
		{"duplicate env files", args{duplicateConfigFiles, ""}, true},
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

func Test_extractEnviron(t *testing.T) {
	env := "environ.local"
	tests := []struct {
		name string
		f    string
		want string
	}{
		{"empty", "", ""},
		{"yml", env + ".yml", env},
		{"yaml", env + ".yaml", env},
		{"json", env + ".json", env},
		{"toml", env + ".toml", env},
		{"dotenv", ".env." + env, env},
		{"noext", env, ""},
		{"invalidext", env + ".bla", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractEnviron(tt.f))
		})
	}
}
