package writers_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thlcodes/genfig/writers"
)

func Test_WriteInit(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, writers.WriteInit(buf, map[string]string{
		"99_A": "A()",
		"20_B": "B()",
		"40_C": "C()",
		"D":    "D()",
	}))
	assert.NotEmpty(t, buf.String())
}
