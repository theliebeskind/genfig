package writers

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/imdario/mergo"

	"github.com/thlcodes/genfig/models"

	"github.com/thlcodes/genfig/util"
	u "github.com/thlcodes/genfig/util"
)

//WriteConfig writes
func WriteConfig(w io.Writer, s models.SchemaMap, config map[string]interface{}, def map[string]interface{}, env string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = u.RecoverError(r)
			return
		}
	}()

	// first write into buffer, so that no gobbledygook is written
	// into some files, when someting panics
	buf := bytes.NewBuffer([]byte{})

	// assigns this config to the according child of the global 'Envs' var
	// via an init function
	buf.Write(u.B("func init() {" + nl + indent + "Envs." + strings.Title(env) + " = "))

	// write actual config
	WriteConfigValue(buf, defaultSchemaRootName, def, config, s, 1)

	merged := def
	if err := mergo.Merge(&merged, config, mergo.WithOverride); err != nil {
		panic(err)
	}
	buf.Write(u.B(nl + indent + "Envs." + strings.Title(env) + "._map = " + fmt.Sprintf("%#v", merged)))

	// closing bracket of init func
	buf.Write(u.B(nl + "}" + nl))

	// now write buffer to writer
	w.Write(buf.Bytes())

	return
}

//WriteConfigLine writes
func WriteConfigLine(w io.Writer, p string, k string, v interface{}, o interface{}, s models.SchemaMap, l int) {
	if l > maxLevel {
		panic(fmt.Errorf("Maximum of %d levels exceeded", maxLevel))
	}

	n := strings.Title(k)

	if _, ex := s[p+n]; !ex {
		panic(fmt.Errorf("Config property '%s' is not defined in the default config", p+n))
	}

	w.Write(u.B(indents[:l*len(indent)]))
	w.Write(u.B(n + ": "))

	WriteConfigValue(w, p+n, v, o, s, l)

	w.Write(u.B("," + nl))
}

//WriteConfigValue writes
func WriteConfigValue(w io.Writer, p string, v interface{}, o interface{}, s models.SchemaMap, l int) {
	switch v.(type) {
	case map[string]interface{}:
		w.Write(u.B(p + "{" + nl))
		for _k, _v := range v.(map[string]interface{}) {
			_o := getOverwriteEntry(o, _k)
			WriteConfigLine(w, p, _k, _v, _o, s, l+1)
		}
		w.Write(u.B(indents[:l*len(indent)]))
		w.Write(u.B("}"))
	case map[interface{}]interface{}:
		w.Write(u.B(p + "{" + nl))
		for _k, _v := range v.(map[interface{}]interface{}) {
			_o := getOverwriteEntry(o, _k)
			WriteConfigLine(w, p, fmt.Sprintf("%v", _k), _v, _o, s, l+1)
		}
		w.Write(u.B(indents[:l*len(indent)]))
		w.Write(u.B("}"))
	case []interface{}:
		t := &v
		if o != nil {
			t = &o
		}
		typ := util.DetectSliceTypeString((*t).([]interface{}))
		w.Write(u.B(strings.Replace(fmt.Sprintf("%#v", *t), "[]interface {}", typ, 1)))
	default:
		t := &v
		if o != nil {
			t = &o
		}
		fmt.Fprintf(w, `%#v`, *t)
	}
}

func getOverwriteEntry(o interface{}, k interface{}) (r interface{}) {
	if o == nil {
		return
	}
	switch o.(type) {
	case map[string]interface{}:
		r = o.(map[string]interface{})[fmt.Sprintf("%v", k)]
	case map[interface{}]interface{}:
		r = o.(map[interface{}]interface{})[k]
	}
	return
}
