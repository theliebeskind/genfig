package plugins

import (
	"io"
	"text/template"

	"github.com/thlcodes/genfig/models"
)

type configTestPlugin struct {
	tpl *template.Template
}

var (
	configTest = configTestPlugin{
		tpl: template.Must(template.
			New("configTest").
			Parse(`import "testing"

func Test_CurrentConfig(t *testing.T) {
	if Current == nil {
		t.Error("Current config is nil")
	}
}
`))}
)

func init() {
	// "register" plugin
	Plugins["90_config_test"] = &configTest
}

// GetInitCall returns the availibility and the string of the
// function to be called on init
func (p *configTestPlugin) GetInitCall() (string, bool) {
	return "", false
}

// SetSchemaMap sets the schema to be used when WriteTo is called
func (p *configTestPlugin) SetSchemaMap(_ models.SchemaMap) {
}

// WriteTo performs the acutal writing to a buffer (or io.Writer).
// For this plugin, the template is simply "rendered" into the writer.
func (p *configTestPlugin) WriteTo(w io.Writer) (l int64, err error) {
	err = p.tpl.Execute(w, nil)
	return
}
