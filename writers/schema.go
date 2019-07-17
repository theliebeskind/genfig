package writers

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/theliebeskind/genfig/models"
	u "github.com/theliebeskind/genfig/util"
)

const (
	defaultSchemaRootName = "Config"
)

//WriteAndReturnSchema writes
func WriteAndReturnSchema(w io.Writer, c map[string]interface{}) (s models.SchemaMap, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = u.RecoverError(r)
			return
		}
	}()

	s = models.SchemaMap{}
	// using NoopWriter since top level is not needed
	WriteSchema(u.NoopWriter{}, defaultSchemaRootName, c, s, 0)

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
			buf.Write(u.B("type " + strings.Replace(k, "_", "", -1) + " " + v.Content + nl))
		}
	}

	// now write buffer to writer
	w.Write(buf.Bytes())
	return
}

//WriteSchema writes
func WriteSchema(w io.Writer, k string, v interface{}, s models.SchemaMap, l int) bool {
	if l > maxLevel {
		panic(fmt.Errorf("Maximum of %d levels exceeded", maxLevel))
	}
	b := bytes.NewBuffer([]byte{})
	n := strings.Title(k)
	isStruct := WriteSchemaType(b, n, v, s, l)

	n = strings.Replace(n, "_", "", -1)
	s[n] = models.Schema{
		IsStruct: isStruct,
		Content:  b.String(),
		Path:     k,
	}
	w.Write(b.Bytes())

	return isStruct
}

// WriteSchemaType the type text to a writer and returns, if type is a struct or not
//WriteSchemaType writes
func WriteSchemaType(w io.Writer, p string, v interface{}, s models.SchemaMap, l int) (isStruct bool) {
	switch v.(type) {
	case map[string]interface{}:
		isStruct = true
		buf := bytes.NewBuffer([]byte{})
		w.Write(u.B("struct {" + nl))
		for _k, _v := range v.(map[string]interface{}) {
			_k = strings.Title(_k)
			_isStruct := WriteSchema(buf, p+"_"+_k, _v, s, l+1)
			if _isStruct {
				w.Write(u.B(indent + _k + " " + strings.Replace(p, "_", "", -1) + _k + nl))
			} else {
				w.Write(u.B(indent + _k + " " + buf.String() + nl))
			}
			buf.Reset()
		}
		w.Write(u.B("}" + nl))
	case map[interface{}]interface{}:
		isStruct = true
		buf := bytes.NewBuffer([]byte{})
		w.Write(u.B("struct {" + nl))
		for _k, _v := range v.(map[interface{}]interface{}) {
			_K := strings.Title(fmt.Sprintf("%v", _k))
			_isStruct := WriteSchema(buf, p+"_"+_K, _v, s, l+1)
			if _isStruct {
				w.Write(u.B(indent + _K + " " + strings.Replace(p, "_", "", -1) + _K + nl))
			} else {
				w.Write(u.B(indent + _K + " " + buf.String() + nl))
			}
			buf.Reset()
		}
		w.Write(u.B("}" + nl))
	case []interface{}:
		w.Write(u.B(u.Make64(u.DetectSliceTypeString(v.([]interface{})))))
	default:
		w.Write(u.B(u.Make64(fmt.Sprintf("%T", v))))
	}
	return
}
