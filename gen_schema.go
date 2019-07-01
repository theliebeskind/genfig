package genfig

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/theliebeskind/genfig/util"
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
	buf := bytes.NewBuffer([]byte{})

	writeSchema(buf, defaultSchemaRootName, c, s, 0)

	buf.Reset()
	buf.Write(u.B("type " + defaultSchemaRootName + " " + s[defaultSchemaRootName].Content + nl))
	for k, v := range s {
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

// writeSchemaType the type text to a writer and returns, if type is a struct or not,
// and an error, if one occures
func writeSchemaType(w io.Writer, p string, v interface{}, s SchemaMap, l int) bool {
	kind := reflect.TypeOf(v).Kind()
	typeStr := reflect.TypeOf(v).String()

	isStruct := false

	if kind == reflect.Map {
		isStruct = true
		w.Write(u.B("struct {"))
		rv := reflect.ValueOf(v)
		if rv.Len() > 0 {
			w.Write(u.B(nl))
		}
		iter := rv.MapRange()
		for iter.Next() {
			k := strings.Title(fmt.Sprintf("%v", iter.Key()))
			b := bytes.NewBuffer([]byte{})
			isStruct := writeSchema(b, p+k, iter.Value().Interface(), s, l+1)
			if isStruct {
				w.Write(u.B(defaultIndent + k + " " + p + k + nl))
			} else {
				w.Write(u.B(defaultIndent + k + " " + b.String() + nl))
			}
		}
		w.Write(u.B("}"))
	} else if kind == reflect.Slice && u.IsInterfaceSlice(v) {
		w.Write(u.B(util.DetectSliceTypeString(v.([]interface{}))))
	} else {
		w.Write(u.B(u.Make64(typeStr)))
	}

	return isStruct
}
