package writers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thclodes/genfig/util"
	"github.com/thclodes/genfig/writers"
)

func Test_WriteHeader(t *testing.T) {
	assert.NoError(t, writers.WriteHeader(util.NoopWriter{}, "", ""))
}
