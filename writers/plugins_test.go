package writers_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thclodes/genfig/models"
	"github.com/thclodes/genfig/plugins"
	"github.com/thclodes/genfig/writers"
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
