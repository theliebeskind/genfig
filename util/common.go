package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	zglob "github.com/mattn/go-zglob"
	yaml "gopkg.in/yaml.v2"
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

// ReverseStrings reverses an
func ReverseStrings(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

// ParseString into an interface
func ParseString(s string) interface{} {
	if iv, err := strconv.ParseInt(s, 10, 0); err == nil {
		return iv
	} else if bv, err := strconv.ParseBool(s); err == nil {
		return bv
	} else if ia, ok := ParseArrayString(s); ok {
		return ia
	} else {
		return s
	}
}

// ParseArrayString parses a string representing an array
func ParseArrayString(s string) ([]interface{}, bool) {
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		return nil, false
	}
	r := []interface{}{}
	if err := yaml.Unmarshal([]byte(s), &r); err != nil {
		return nil, false
	}
	return r, true
}

// RecoverError recovers errors
func RecoverError(r interface{}) error {
	switch r.(type) {
	case error:
		return r.(error)
	default:
		return fmt.Errorf("%v", r)
	}
}
