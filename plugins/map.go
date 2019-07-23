package plugins

import (
	"io"
	"text/template"

	"github.com/theliebeskind/genfig/models"
)

type mapPlugin struct {
	s   models.SchemaMap
	tpl *template.Template
}

var (
	mapt = mapPlugin{
		tpl: template.Must(template.
			New("map").
			Parse(`func (c *Config) Map() map[string]interface{} {
	return c._map
}
`))}
)

func init() {
	// "register" plugin
	Plugins["map"] = &mapt
}

// GetInitCall returns the availibility and the string of the
// function to be called on init
func (p *mapPlugin) GetInitCall() (string, bool) {
	return "", false
}

// SetSchemaMap sets the schema to be used when WriteTo is called
func (p *mapPlugin) SetSchemaMap(s models.SchemaMap) {
	p.s = s
}

// WriteTo performs the acutal writing to a buffer (or io.Writer).
// For this plugin, the template is simply "rendered" into the writer.
func (p *mapPlugin) WriteTo(w io.Writer) (l int64, err error) {
	err = p.tpl.Execute(w, nil)
	return
}
