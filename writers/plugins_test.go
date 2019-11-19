package writers_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thlcodes/genfig/models"
	"github.com/thlcodes/genfig/plugins"
	"github.com/thlcodes/genfig/writers"
)

func Test_WritePlugins(t *testing.T) {
	sm := models.SchemaMap{}
	calls := map[string]string{}
	dir, _ := ioutil.TempDir("", "genfig")
	defer os.RemoveAll(dir)
	files, err := writers.WritePlugins(sm, dir, "test", "genfig test", calls)
	assert.NoError(t, err)
	assert.Len(t, files, len(plugins.Plugins))
	for _, f := range files {
		assert.FileExists(t, f)
	}
}
