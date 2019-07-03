package types

import "io"

// Schema defines the schema
type Schema struct {
	IsStruct bool
	Content  string
	Path     string
}

// SchemaMap aliases as string-map of bytes
type SchemaMap map[string]Schema

// Params for the Generate func as struct,
// empty values are default values, so can be passed empty
type Params struct {
	Dir        string
	DefaultEnv string
	MergeFiles bool
}

// Plugin interface
type Plugin interface {
	io.WriterTo
	SetSchemaMap(SchemaMap)
}
