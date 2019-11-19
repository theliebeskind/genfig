package writers_test

import (
	"strings"
	"testing"

	"github.com/thlcodes/genfig/writers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thlcodes/genfig/models"
	"github.com/thlcodes/genfig/util"
)

func Test_WriteConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		def      map[string]interface{}
		contains []string
		wantErr  bool
	}{
		{"empty", map[string]interface{}{}, nil, []string{}, false},
		{"simple string", map[string]interface{}{"a": "b"}, nil, []string{"A: \"b\""}, false},
		{"simple string with default", map[string]interface{}{"a": "b"}, map[string]interface{}{"a": "def"}, []string{"A: \"b\""}, false},
		{"simple int", map[string]interface{}{"a": 1}, nil, []string{"A: 1"}, false},
		{"simple bool", map[string]interface{}{"a": true}, nil, []string{"A: true"}, false},
		{"int array", map[string]interface{}{"a": []int{1, 2, 3}}, nil, []string{"A: []int"}, false},
		{"empy interface array", map[string]interface{}{"a": []interface{}{}}, nil, []string{"A: []interface {}{}"}, false},
		{"mixed interface array", map[string]interface{}{"a": []interface{}{1, ""}}, nil, []string{"A: []interface {}{1, \"\"}"}, false},
		{"int interface array", map[string]interface{}{"a": []interface{}{1, 2}}, nil, []string{"A: []int{1, 2}"}, false},
		{"string interface array", map[string]interface{}{"a": []interface{}{"a", "b"}}, nil, []string{"A: []string"}, false},
		{"map", map[string]interface{}{"a": map[string]interface{}{"b": 1}}, nil, []string{"A: ConfigA{", "B: 1"}, false},
		{"map with interface key", map[string]interface{}{"a": map[interface{}]interface{}{"b": 1}}, nil, []string{"A: ConfigA{", "B: 1"}, false},
		{"map of map", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, nil, []string{"A: ConfigA{", "B: ConfigAB{", "C: 1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &strings.Builder{}
			def := tt.def
			if def == nil {
				def = tt.config
			}
			err := writers.WriteConfig(s, models.SchemaMap{
				"ConfigA":   models.Schema{},
				"ConfigAB":  models.Schema{},
				"ConfigABC": models.Schema{},
			}, tt.config, def, "test")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			got := s.String()
			require.NoError(t, err)
			for _, c := range tt.contains {
				assert.Contains(t, got, c)
			}
		})
	}
}

func Benchmark_WriteConfigValue(b *testing.B) {
	w := util.NoopWriter{}
	s := models.SchemaMap{
		"ConfigA":   models.Schema{},
		"ConfigAB":  models.Schema{},
		"ConfigABC": models.Schema{},
		"ConfigABD": models.Schema{},
		"ConfigABE": models.Schema{},
	}
	m := map[string]interface{}{"a": map[interface{}]interface{}{"b": map[string]interface{}{"c": []interface{}{1}, "d": "s", "e": 1}}}
	e := map[string]interface{}{}
	for n := 0; n < b.N; n++ {
		writers.WriteConfigValue(w, "Config", m, e, s, 0)
	}
}
