package writers

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/imdario/mergo"

	"github.com/thlcodes/genfig/models"

	"github.com/thlcodes/genfig/util"
	u "github.com/thlcodes/genfig/util"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}

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

	merged := make(map[string]interface{})
	copyMap(def, &merged)
	if err := mergo.Merge(&merged, config, mergo.WithOverride); err != nil {
		panic(err)
	}

	// write actual config
	WriteConfigValue(buf, defaultSchemaRootName, merged, s, 1)

	buf.Write(u.B(nl + indent + "Envs." + strings.Title(env) + "._map = " + fmt.Sprintf("%#v", merged)))

	// closing bracket of init func
	buf.Write(u.B(nl + "}" + nl))

	// now write buffer to writer
	w.Write(buf.Bytes())

	return
}

//WriteConfigLine writes
func WriteConfigLine(w io.Writer, p string, k string, v interface{}, s models.SchemaMap, l int) {
	if l > maxLevel {
		panic(fmt.Errorf("Maximum of %d levels exceeded", maxLevel))
	}

	n := strings.Title(k)

	if _, ex := s[p+n]; !ex {
		panic(fmt.Errorf("Config property '%s' is not defined in the default config", p+n))
	}

	w.Write(u.B(indents[:l*len(indent)]))
	w.Write(u.B(n + ": "))

	WriteConfigValue(w, p+n, v, s, l)

	w.Write(u.B("," + nl))
}

//WriteConfigValue writes
func WriteConfigValue(w io.Writer, p string, v interface{}, s models.SchemaMap, l int) {
	switch v.(type) {
	case map[string]interface{}:
		w.Write(u.B(p + "{" + nl))
		keys := []string{}
		for _k := range v.(map[string]interface{}) {
			keys = append(keys, _k)
		}
		sort.Strings(keys)
		for _, _k := range keys {
			_v := v.(map[string]interface{})[_k]
			//_o := getOverwriteEntry(o, _k)
			WriteConfigLine(w, p, _k, _v /*, _o*/, s, l+1)
		}
		w.Write(u.B(indents[:l*len(indent)]))
		w.Write(u.B("}"))
	case []interface{}:
		t := &v
		typ := util.DetectSliceTypeString((*t).([]interface{}))
		w.Write(u.B(strings.Replace(fmt.Sprintf("%#v", *t), "[]interface {}", typ, 1)))
	default:
		t := &v
		fmt.Fprintf(w, `%#v`, *t)
	}
}

func copyMap(from map[string]interface{}, to *map[string]interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(from)
	if err != nil {
		return err
	}
	err = dec.Decode(to)
	if err != nil {
		return err
	}
	return nil
}
