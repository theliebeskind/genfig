package writers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/thlcodes/genfig/plugins"

	"github.com/thlcodes/genfig/models"
)

const (
	pluginPrefix = "plugin_"
)

//WritePlugins writes a plugin file for each plugin
func WritePlugins(schema models.SchemaMap, dir string, pkg string, cmd string, calls map[string]string) ([]string, error) {
	files := []string{}
	for n, p := range plugins.Plugins {
		orig := n
		p.SetSchemaMap(schema)
		if strings.Contains(n, "_") {
			n = n[strings.Index(n, "_")+1:]
		}
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
		if c, has := p.GetInitCall(); has {
			calls[orig] = c
		}
	}
	return files, nil
}
