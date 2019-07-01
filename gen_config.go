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

func writeConfig(w io.Writer, s SchemaMap, c map[string]interface{}, env string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = u.RecoverError(r)
			return
		}
	}()

	buf := bytes.NewBuffer([]byte{})

	buf.Write(u.B("var " + strings.Title(env) + " = "))

	writeConfigValue(buf, defaultSchemaRootName, c, s, 0)

	// now write buffer to writer
	w.Write(buf.Bytes())
	return
}

func writeConfigLine(w io.Writer, p string, k string, v interface{}, s SchemaMap, l int) {
	if l > maxLevel {
		panic(fmt.Errorf("Maximum of %d levels exceeded", maxLevel))
	}

	n := strings.Title(k)

	if _, ex := s[p+n]; !ex {
		panic(fmt.Errorf("Config property '%s' is not defined in the default config", p+n))
	}

	w.Write(u.B(strings.Repeat(defaultIndent, l)))
	w.Write(u.B(n + ": "))

	writeConfigValue(w, p+n, v, s, l)

	w.Write(u.B("," + nl))
}

// writeConfigValue writes the actual value
func writeConfigValue(w io.Writer, p string, v interface{}, s SchemaMap, l int) {
	kind := reflect.TypeOf(v).Kind()

	if kind == reflect.Map {
		w.Write(u.B(p + "{" + nl))
		rv := reflect.ValueOf(v)
		if rv.Len() > 0 {
			iter := rv.MapRange()
			for iter.Next() {
				k := strings.Title(fmt.Sprintf("%v", iter.Key()))
				writeConfigLine(w, p, k, iter.Value().Interface(), s, l+1)
			}
		}
		w.Write(u.B(strings.Repeat(defaultIndent, l)))
		w.Write(u.B("}"))
	} else if kind == reflect.Slice && u.IsInterfaceSlice(v) {
		typ := util.DetectSliceTypeString(v.([]interface{}))
		w.Write(u.B(strings.Replace(fmt.Sprintf("%#v", v), "[]interface {}", typ, 1)))
	} else if kind == reflect.Slice {
		fmt.Fprintf(w, "%#v", v)
	} else if kind == reflect.String {
		fmt.Fprintf(w, `"%v"`, v)
	} else {
		fmt.Fprintf(w, "%v", v)
	}
}
