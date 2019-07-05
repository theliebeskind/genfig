package plugins

import (
	"io"

	"github.com/theliebeskind/genfig/types"
)

// Plugin interface
type Plugin interface {
	io.WriterTo
	SetSchemaMap(types.SchemaMap)
	GetInitCall() (string, bool)
}

// Plugins hold the available plugins
var Plugins = map[string]Plugin{}
