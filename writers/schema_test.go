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

func Test_WriteSchema(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]interface{}
		contains   []string
		wantSchema models.Schema
		wantErr    bool
	}{
		{"empty", map[string]interface{}{}, []string{}, models.Schema{}, false},
		{"simple string", map[string]interface{}{"a": "b"}, []string{"A string"}, models.Schema{}, false},
		{"simple int", map[string]interface{}{"a": 1}, []string{"A int64"}, models.Schema{}, false},
		{"simple bool", map[string]interface{}{"a": true}, []string{"A bool"}, models.Schema{}, false},
		{"int array", map[string]interface{}{"a": []int{1, 2, 3}}, []string{"A []int"}, models.Schema{}, false},
		{"empty interface array", map[string]interface{}{"a": []interface{}{}}, []string{"A []interface {}"}, models.Schema{}, false},
		{"mixed interface array", map[string]interface{}{"a": []interface{}{1, ""}}, []string{"A []interface {}"}, models.Schema{}, false},
		{"int interface array", map[string]interface{}{"a": []interface{}{1, 2}}, []string{"A []int64"}, models.Schema{}, false},
		{"string interface array", map[string]interface{}{"a": []interface{}{"a", "b"}}, []string{"A []string"}, models.Schema{}, false},
		{"map", map[string]interface{}{"a": map[string]interface{}{"b": 1}}, []string{"A struct {", "B int"}, models.Schema{}, false},
		{"iface key map", map[string]interface{}{"a": map[interface{}]interface{}{"b": 1}}, []string{"A struct {", "B int"}, models.Schema{}, false},
		{"map of map", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, []string{"A struct {", "B struct {", "C int"}, models.Schema{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &strings.Builder{}
			_, err := writers.WriteAndReturnSchema(s, tt.config)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			got := s.String()
			for _, c := range tt.contains {
				assert.Contains(t, got, c)
			}
		})
	}
}

func Benchmark_WriteSchemaType(b *testing.B) {
	w := util.NoopWriter{}
	s := models.SchemaMap{}
	m := map[string]interface{}{"a": map[interface{}]interface{}{"b0": 1, "b": map[string]interface{}{"c": []interface{}{1}, "d": "s", "e": 1}}}
	for n := 0; n < b.N; n++ {
		writers.WriteSchemaType(w, "Config", m, s, 0)
	}
}
