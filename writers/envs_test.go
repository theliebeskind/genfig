package writers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thlcodes/genfig/util"
	"github.com/thlcodes/genfig/writers"
)

func Test_WriteEnvs(t *testing.T) {
	assert.NoError(t, writers.WriteEnvs(util.NoopWriter{}, map[string]string{}))
}
