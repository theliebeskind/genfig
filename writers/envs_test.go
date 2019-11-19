package writers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thclodes/genfig/util"
	"github.com/thclodes/genfig/writers"
)

func Test_WriteEnvs(t *testing.T) {
	assert.NoError(t, writers.WriteEnvs(util.NoopWriter{}, map[string]string{}))
}
