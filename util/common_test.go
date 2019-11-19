package util_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thlcodes/genfig/util"
)

var (
	fixturesDir, _ = filepath.Abs("../fixtures")
)

func Test_ResolveGlobs(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "genifig")
	defer os.RemoveAll(tmpDir)

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	_ = os.Chdir(tmpDir)

	ioutil.WriteFile("a.x", []byte{}, 0777)
	ioutil.WriteFile("b.x", []byte{}, 0777)
	ioutil.WriteFile("c.y", []byte{}, 0777)
	os.MkdirAll("sub", 0777)
	ioutil.WriteFile("sub/d.y", []byte{}, 0777)

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
		{"multiple globs, unique", args{[]string{"a.*", "b.*"}}, []string{"a.x", "b.x"}},
		{"multiple globs, not unique", args{[]string{"a.*", "*.x"}}, []string{"a.x", "b.x"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.ResolveGlobs(tt.args.globs...)
			sort.Strings(tt.want)
			sort.Strings(got)
			assert.Equal(t, tt.want, got)
		})
	}
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

func Test_CleanDir(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "genfig")
	defer os.RemoveAll(tmpDir)
	dir := "clean"
	cwd, _ := os.Getwd()
	os.MkdirAll(filepath.Join(tmpDir, dir), 0777)
	os.Chdir(filepath.Join(tmpDir, dir))
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

func Test_ParseString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want interface{}
	}{
		{"empty", "", ""},
		{"int", "1", int64(1)},
		{"float", "1.1", float64(1.1)},
		{"negative int", "-999", int64(-999)},
		{"bool true", "true", true},
		{"bool false", "false", false},
		{"int array", "[1,2,3]", []interface{}{float64(1), float64(2), float64(3)}},
		{"string array", "[\"a\",\"b\",\"c\"]", []interface{}{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.ParseString(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ParseStringArray(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		want   []interface{}
		wantOk bool
	}{
		{"empty data", "", nil, false},
		{"not array", "a,b", nil, false},
		{"invalid array", "[a,sä#.,,]", nil, false},
		{"empty array", "[]", []interface{}{}, true},
		{"ints", "[1,2,3]", []interface{}{float64(1), float64(2), float64(3)}, true},
		{"floats", "[1.1,2.2,3.3]", []interface{}{1.1, 2.2, 3.3}, true},
		{"bools", "[true, false, true]", []interface{}{true, false, true}, true},
		{"strings", `["a", "b", "c"]`, []interface{}{"a", "b", "c"}, true},
		{"mixed", "[true, 1, 2.2, \"s\"]", []interface{}{true, float64(1), 2.2, "s"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := util.ParseArrayString(tt.s)
			if tt.wantOk {
				require.True(t, ok)
			} else {
				require.False(t, ok)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Reverse(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", []string{}, []string{}},
		{"even", []string{"a", "b", "c", "d"}, []string{"d", "c", "b", "a"}},
		{"odd", []string{"a", "b", "c"}, []string{"c", "b", "a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			util.ReverseStrings(tt.in)
			assert.Equal(t, tt.want, tt.in)
		})
	}
}

func Test_DetectSliceTypeString(t *testing.T) {
	tests := []struct {
		name  string
		slice []interface{}
		want  string
	}{
		{"empty", []interface{}{}, "[]interface {}"},
		{"ints", []interface{}{1, 2, 3}, "[]int"},
		{"bools", []interface{}{true, false, false}, "[]bool"},
		{"string", []interface{}{"a", "b", ""}, "[]string"},
		{"mixed", []interface{}{"a", 1, false}, "[]interface {}"},
		{"structs", []interface{}{struct{ a int }{}}, "[]struct { a int }"},
		{"maps", []interface{}{map[string]interface{}{}}, "[]map[string]interface {}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.DetectSliceTypeString(tt.slice)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_IsInterfaceSlice(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		want bool
	}{
		{"no", "nope", false},
		{"also no", struct{}{}, false},
		{"yes", []interface{}{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, util.IsInterfaceSlice(tt.in))
		})
	}
}

func Test_Make64(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"int", "int64"},
		{"uint", "uint64"},
		{"float", "float64"},
		{"[]float", "[]float64"},
		{"map[string]int", "map[string]int64"},
		{"string", "string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, util.Make64(tt.name))
		})
	}
}
