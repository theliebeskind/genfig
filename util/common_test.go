package util_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theliebeskind/genfig/util"
)

const (
	fixturesDir = "../fixtures/"
)

func Test_ResolveGlobs(t *testing.T) {
	type args struct {
		globs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"no args", args{}, []string{}},
		{"empty arg", args{[]string{""}}, []string{}},
		{"no match", args{[]string{"foo"}}, []string{}},
		{"all first level", args{[]string{"*"}}, []string{"c.y", "a.x", "b.x"}},
		{"y first level file", args{[]string{"*.y"}}, []string{"c.y"}},
		{"x first leve files", args{[]string{"*.x"}}, []string{"a.x", "b.x"}},
		{"all y files", args{[]string{"**/*.y"}}, []string{"c.y", "sub/d.y"}},
		{"multipe globs, unique", args{[]string{"a.*", "b.*"}}, []string{"a.x", "b.x"}},
		{"multipe globs, not unique", args{[]string{"a.*", "*.x"}}, []string{"a.x", "b.x"}},
	}
	cwd, _ := os.Getwd()
	_ = os.Chdir(filepath.Join(fixturesDir, "/dir"))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.ResolveGlobs(tt.args.globs...)
			sort.Strings(tt.want)
			sort.Strings(got)
			assert.Equal(t, tt.want, got)
		})
	}
	_ = os.Chdir(cwd)
}

func Test_MapString(t *testing.T) {
	type args struct {
		vs []string
		f  func(string) string
	}
	f := func(s string) string {
		return s + "_"
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"nil", args{nil, f}, []string{}},
		{"empty", args{[]string{}, f}, []string{}},
		{"normal", args{[]string{"a", "b", "_", ""}, f}, []string{"a_", "b_", "__", "_"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, util.MapString(tt.args.vs, tt.args.f))
		})
	}
}

func Test_ReduceStrings(t *testing.T) {
	type args struct {
		vs []string
		f  func(interface{}, string) interface{}
		r  interface{}
	}
	f1 := func(r interface{}, i string) interface{} {
		return r.(string) + i
	}
	f2 := func(r interface{}, i string) interface{} {
		if i != "" {
			r = append(r.([]string), i)
		}
		return r
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"nil", args{nil, f1, ""}, ""},
		{"empty", args{[]string{}, f1, ""}, ""},
		{"string", args{[]string{"a", "b", "_", ""}, f1, ""}, "ab_"},
		{"slice", args{[]string{"a", "b", "_", ""}, f2, []string{}}, []string{"a", "b", "_"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, util.ReduceStrings(tt.args.vs, tt.args.f, tt.args.r))
		})
	}
}

func TestCleanDir(t *testing.T) {
	dir := "clean"
	cwd, _ := os.Getwd()
	os.Chdir(filepath.Join(fixturesDir, dir))
	defer os.Chdir(cwd)

	ioutil.WriteFile("a.txt", []byte{}, 0777)
	ioutil.WriteFile("b.txt", []byte{}, 0777)
	assert.NotEmpty(t, util.ResolveGlobs("*"))

	os.Chdir("..")

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nonexist", args{"nonexist"}, true},
		{"exist", args{dir}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.CleanDir(tt.args.name)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.Empty(t, util.ResolveGlobs("./"+tt.args.name+"/*"))
		})
	}
}
