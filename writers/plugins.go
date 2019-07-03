package writers

import (
	"os"
	"path/filepath"

	"github.com/theliebeskind/genfig/plugins"

	"github.com/theliebeskind/genfig/types"
)

const (
	pluginPrefix = "plugin_"
)

//WritePlugins writes
func WritePlugins(schema types.SchemaMap, dir string, pkg string, cmd string) ([]string, error) {
	files := []string{}
	for n, p := range plugins.Plugins {
		p.SetSchemaMap(schema)
		path := filepath.Join(dir, pluginPrefix+n+".go")
		if f, err := os.Create(path); err != nil {
			return files, err
		} else if err := WriteHeader(f, pkg, cmd+" plugin '"+n+"'"); err != nil {
			return files, err
		} else if _, err := p.WriteTo(f); err != nil {
			return files, err
		} else {
			_ = f.Close()
			files = append(files, path)
		}
	}
	return files, nil
}
