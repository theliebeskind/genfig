package writers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thlcodes/genfig/util"
	"github.com/thlcodes/genfig/writers"
)

func Test_WriteHeader(t *testing.T) {
	assert.NoError(t, writers.WriteHeader(util.NoopWriter{}, "", ""))
}
