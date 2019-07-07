package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

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
	} else if fv, err := strconv.ParseFloat(s, 0); err == nil {
		return fv
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
	if err := json.Unmarshal([]byte(s), &r); err != nil {
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

// DetectSliceTypeString returns the actual type of an slice of interfaces
func DetectSliceTypeString(slice []interface{}) string {
	iface := "[]interface {}"
	if len(slice) == 0 {
		return iface
	}
	var typ reflect.Type
	for _, s := range slice {
		t := reflect.TypeOf(s)
		if typ == nil {
			typ = t
			continue
		}
		if t != typ {
			return iface
		}
	}
	return "[]" + typ.String()
}

// IsInterfaceSlice checks if a given interface is actually a slice of interfaces
func IsInterfaceSlice(i interface{}) (is bool) {
	_, is = i.([]interface{})
	return
}

// Make64 adds '64' to 'int', 'uint' and 'float'
func Make64(s string) string {
	if strings.HasSuffix(s, "int") || strings.HasSuffix(s, "uint") || strings.HasSuffix(s, "float") {
		return s + "64"
	}
	return s
}

// B is an unsafe string to bytes
func B(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// NoopWriter is a writer that does nothing
type NoopWriter struct{}

func (NoopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
