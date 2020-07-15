package plugins

import (
	"io"
	"regexp"
	"strings"
	"text/template"

	"github.com/thlcodes/genfig/models"
)

type updateFromEnvPlugin struct {
	s   models.SchemaMap
	tpl *template.Template
}

var (
	sliceMatcher = regexp.MustCompile(`^\[\](\w+)`)
)

var (
	updateFromEnv = updateFromEnvPlugin{
		s: models.SchemaMap{},
		tpl: template.Must(template.
			New("updateFromEnv").
			Funcs(template.FuncMap{
				"upper": strings.ToUpper,
				"lower": strings.ToLower,
				"title": strings.Title,
				// Remove root (usually "Config_") from env var name
				"cleanPrefixEnv": func(s string) string {
					return strings.Join(strings.Split(s, "_")[1:], "_")
				},
				// A_B to a.b
				"dotPath": func(s string) string {
					return strings.ReplaceAll(s, "_", ".")
				},
				// Converte an env var name to a Config path
				"makePath": func(s string) string {
					return strings.Join(strings.Split(s, "_")[1:], ".")
				},
				// Substitute []*type* with *type*Slice
				"renameSlice": func(s string) string {
					if found := sliceMatcher.FindStringSubmatch(s); len(found) > 0 {
						return found[1] + "Slice"
					}
					return s
				},
			}).
			Parse(`import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	_ = os.LookupEnv
	_ = fmt.Sprintf
	_ = json.Marshal
)

func (c *Config) UpdateFromEnv() []error {
	var val string
	_ = val
	var exists bool
	_ = exists
	var envs []string
	_ = envs
	var errors = []error{}
{{range $_, $v := .}}{{if not $v.IsStruct}}
	envs = []string{"{{dotPath (cleanPrefixEnv (lower $v.Path))}}", "{{cleanPrefixEnv (upper $v.Path)}}"}
	for _, env := range envs {
		if val, exists = os.LookupEnv(env); exists {
			break
		}
	}
	if exists { {{if eq $v.Content "string"}}
		c.{{makePath $v.Path}} = val {{else}}
		if err := parse{{title (renameSlice $v.Content)}}(val, &c.{{makePath $v.Path}}); err != nil {
			errors = append(errors, fmt.Errorf("Genfig: could not parse {{$v.Content}} from {{upper $v.Path}} ('%s')\n", val))
		} {{end}}
	}
{{end}}{{end}}
	if len(errors) == 0 {
		return nil
	} else {
		return errors
	}
}

// these are wrappers, so that they can
// a) be referenced easily be the code generator and
// b) be replaces easily by you (or me)
func parseInt64(s string, i *int64) (err error) {
	if got, err := strconv.ParseInt(s, 10, 0); err == nil {
		*i = got
	}
	return
}

func parseFloat64(s string, f *float64) (err error) {
	if got, err := strconv.ParseFloat(s, 0); err == nil {
		*f = got
	}
	return
}

func parseBool(s string, b *bool) (err error) {
	if got, err := strconv.ParseBool(s); err != nil {
		*b = got
	}
	return
}

func parseStringSlice(s string, a *[]string) (err error) {
	add := false
	if strings.HasPrefix(s, "+") {
		add = true
		s = s[1:]
	}
	tmp := []string{}
	if err = json.Unmarshal([]byte(s), &tmp); err != nil {
		return
	}
	if add {
		*a = append(*a, tmp...)
	} else {
		*a = tmp
	}
	return
}

func parseInt64Slice(s string, a *[]int64) (err error) {
	add := false
	if strings.HasPrefix(s, "+") {
		add = true
		s = s[1:]
	}
	tmp := []int64{}
	if err = json.Unmarshal([]byte(s), &tmp); err != nil {
		return
	}
	if add {
		*a = append(*a, tmp...)
	} else {
		*a = tmp
	}
	return
}

func parseFloat64Slice(s string, a *[]float64) (err error) {
	add := false
	if strings.HasPrefix(s, "+") {
		add = true
		s = s[1:]
	}
	tmp := []float64{}
	if err = json.Unmarshal([]byte(s), &tmp); err != nil {
		return
	}
	if add {
		*a = append(*a, tmp...)
	} else {
		*a = tmp
	}
	return
}

func parseInterfaceSlice(s string, a *[]interface{},) (err error) {
	add := false
	if strings.HasPrefix(s, "+") {
		add = true
		s = s[1:]
	}
	tmp := []interface{}{}
	if err = json.Unmarshal([]byte(s), &tmp); err != nil {
		return
	}
	if add {
		*a = append(*a, tmp...)
	} else {
		*a = tmp
	}
	return
}

func parseMapSlice(s string, a *[]map[string]interface{}) (err error) {
	add := false
	if strings.HasPrefix(s, "+") {
		add = true
		s = s[1:]
	}
	tmp := []map[string]interface{}{}
	if err = json.Unmarshal([]byte(s), &tmp); err != nil {
		return
	}
	if add {
		*a = append(*a, tmp...)
	} else {
		*a = tmp
	}
	return
}
`))}
)

func init() {
	// "register" plugin
	Plugins["30_update_from_env"] = &updateFromEnv
}

// GetInitCall returns the availibility and the string of the
// function to be called on init
func (p *updateFromEnvPlugin) GetInitCall() (string, bool) {
	return "if errs := Current.UpdateFromEnv(); errs != nil {\n\tfmt.Println(errs)\n}", true
}

// SetSchemaMap sets the schema to be used when WriteTo is called
func (p *updateFromEnvPlugin) SetSchemaMap(s models.SchemaMap) {
	p.s = s
}

// WriteTo performs the acutal writing to a buffer (or io.Writer).
// For this plugin, the template is simply "rendered" into the writer.
func (p *updateFromEnvPlugin) WriteTo(w io.Writer) (l int64, err error) {
	err = p.tpl.Execute(w, p.s)
	return
}
