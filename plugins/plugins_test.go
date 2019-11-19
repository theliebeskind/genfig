package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thlcodes/genfig/models"
	"github.com/thlcodes/genfig/util"

	"github.com/thlcodes/genfig/plugins"
)

func Test_All(t *testing.T) {
	s := models.SchemaMap{
		"A": models.Schema{
			Content: "string",
			Path:    "A",
		},
		"B": models.Schema{
			Content:  "struct { C int }",
			Path:     "B",
			IsStruct: true,
		},
		"BC": models.Schema{
			Content: "int",
			Path:    "BC",
		},
	}
	for _, p := range plugins.Plugins {
		p.SetSchemaMap(s)
		c, b := p.GetInitCall()
		if b {
			assert.NotEmpty(t, c)
		}
		assert.NotPanics(t, func() {
			_, err := p.WriteTo(util.NoopWriter{})
			assert.NoError(t, err)
		})
	}
}
