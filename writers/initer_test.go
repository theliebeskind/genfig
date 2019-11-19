package writers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thclodes/genfig/util"
	"github.com/thclodes/genfig/writers"
)

func Test_WriteInit(t *testing.T) {
	assert.NoError(t, writers.WriteInit(util.NoopWriter{}, map[string]string{}))
}
