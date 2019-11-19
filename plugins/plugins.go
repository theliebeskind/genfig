package plugins

import (
	"io"

	"github.com/thclodes/genfig/models"
)

// Plugin interface
type Plugin interface {
	io.WriterTo
	SetSchemaMap(models.SchemaMap)
	GetInitCall() (string, bool)
}

// Plugins hold the available plugins
var Plugins = map[string]Plugin{}
