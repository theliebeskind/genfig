package util

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	zglob "github.com/mattn/go-zglob"
)

// ResolveGlobs resolves globs and returns all found files unique
func ResolveGlobs(globs ...string) []string {
	m := map[string]struct{}{}
	for _, glob := range globs {
		found, err := zglob.Glob(glob) // using zglob since filepath.Glob did not work with double star
		if err != nil {
			continue
		}
		for _, f := range found {
			m[f] = struct{}{}
		}
	}
	files := []string{}
	for k := range m {
		if info, _ := os.Stat(k); !info.IsDir() && !strings.HasSuffix(k, ".DS_Store") {
			files = append(files, k)
		}
	}
	return files
}

// MapString maps an array of strings
func MapString(vs []string, f func(s string) string) []string {
	if len(vs) == 0 || vs == nil {
		return make([]string, 0)
	}
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// ReduceStrings reduces a string array to a string
func ReduceStrings(vs []string, f func(r interface{}, s string) interface{}, r interface{}) interface{} {
	if len(vs) == 0 || vs == nil {
		return r
	}
	for i := range vs {
		r = f(r, vs[i])
	}
	return r
}

// CleanDir cleans a directory
func CleanDir(name string) error {
	dir, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}
	for _, d := range dir {
		err := os.RemoveAll(path.Join([]string{name, d.Name()}...))
		if err != nil {
			return err
		}
	}
	return nil
}
