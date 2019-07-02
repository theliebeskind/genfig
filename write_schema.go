package genfig

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	u "github.com/theliebeskind/genfig/util"
)

const (
	defaultSchemaRootName = "Config"
)

func writeAndReturnSchema(w io.Writer, c map[string]interface{}) (s SchemaMap, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = u.RecoverError(r)
			return
		}
	}()

	s = SchemaMap{}
	// using NoopWriter since top level is not needed
	writeSchema(u.NoopWriter{}, defaultSchemaRootName, c, s, 0)

	buf := bytes.NewBuffer([]byte{})
	// write top level schema type definition (usually 'Config')
	buf.Write(u.B("type " + defaultSchemaRootName + " " + s[defaultSchemaRootName].Content + nl))
	keys := []string{}
	for k := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := s[k]
		if v.IsStruct {
			if k == defaultSchemaRootName {
				continue
			}
			buf.Write(u.B("type " + k + " " + v.Content + nl))
		}
	}

	// now write buffer to writer
	w.Write(buf.Bytes())
	return
}

func writeSchema(w io.Writer, k string, v interface{}, s SchemaMap, l int) bool {
	if l > maxLevel {
		panic(fmt.Errorf("Maximum of %d levels exceeded", maxLevel))
	}
	b := bytes.NewBuffer([]byte{})
	n := strings.Title(k)
	isStruct := writeSchemaType(b, n, v, s, l)

	s[n] = Schema{isStruct, b.String()}
	w.Write(b.Bytes())

	return isStruct
}

// writeSchemaType the type text to a writer and returns, if type is a struct or not
func writeSchemaType(w io.Writer, p string, v interface{}, s SchemaMap, l int) (isStruct bool) {
	switch v.(type) {
	case map[string]interface{}:
		isStruct = true
		buf := bytes.NewBuffer([]byte{})
		w.Write(u.B("struct {" + nl))
		for _k, _v := range v.(map[string]interface{}) {
			_k = strings.Title(_k)
			_isStruct := writeSchema(buf, p+_k, _v, s, l+1)
			if _isStruct {
				w.Write(u.B(defaultIndent + _k + " " + p + _k + nl))
			} else {
				w.Write(u.B(defaultIndent + _k + " " + buf.String() + nl))
			}
			buf.Reset()
		}
		w.Write(u.B("}"))
	case map[interface{}]interface{}:
		isStruct = true
		buf := bytes.NewBuffer([]byte{})
		w.Write(u.B("struct {" + nl))
		for _k, _v := range v.(map[interface{}]interface{}) {
			_K := strings.Title(fmt.Sprintf("%v", _k))
			_isStruct := writeSchema(buf, p+_K, _v, s, l+1)
			if _isStruct {
				w.Write(u.B(defaultIndent + _K + " " + p + _K + nl))
			} else {
				w.Write(u.B(defaultIndent + _K + " " + buf.String() + nl))
			}
			buf.Reset()
		}
		w.Write(u.B("}"))
	case []interface{}:
		w.Write(u.B(u.Make64(u.DetectSliceTypeString(v.([]interface{})))))
	default:
		// get type string via reflect
		w.Write(u.B(u.Make64(fmt.Sprintf("%T", v))))
	}
	return
}
