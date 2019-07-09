package writers_test

import (
	"strings"
	"testing"

	"github.com/theliebeskind/genfig/writers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theliebeskind/genfig/types"
	"github.com/theliebeskind/genfig/util"
)

func Test_WriteSchema(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]interface{}
		contains   []string
		wantSchema types.Schema
		wantErr    bool
	}{
		{"empty", map[string]interface{}{}, []string{}, types.Schema{}, false},
		{"simple string", map[string]interface{}{"a": "b"}, []string{"A string"}, types.Schema{}, false},
		{"simple int", map[string]interface{}{"a": 1}, []string{"A int64"}, types.Schema{}, false},
		{"simple bool", map[string]interface{}{"a": true}, []string{"A bool"}, types.Schema{}, false},
		{"int array", map[string]interface{}{"a": []int{1, 2, 3}}, []string{"A []int"}, types.Schema{}, false},
		{"empty interface array", map[string]interface{}{"a": []interface{}{}}, []string{"A []interface {}"}, types.Schema{}, false},
		{"mixed interface array", map[string]interface{}{"a": []interface{}{1, ""}}, []string{"A []interface {}"}, types.Schema{}, false},
		{"int interface array", map[string]interface{}{"a": []interface{}{1, 2}}, []string{"A []int64"}, types.Schema{}, false},
		{"string interface array", map[string]interface{}{"a": []interface{}{"a", "b"}}, []string{"A []string"}, types.Schema{}, false},
		{"map", map[string]interface{}{"a": map[string]interface{}{"b": 1}}, []string{"A struct {", "B int"}, types.Schema{}, false},
		{"iface key map", map[string]interface{}{"a": map[interface{}]interface{}{"b": 1}}, []string{"A struct {", "B int"}, types.Schema{}, false},
		{"map of map", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": 1}}}, []string{"A struct {", "B struct {", "C int"}, types.Schema{}, false},
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
	s := types.SchemaMap{}
	m := map[string]interface{}{"a": map[interface{}]interface{}{"b0": 1, "b": map[string]interface{}{"c": []interface{}{1}, "d": "s", "e": 1}}}
	for n := 0; n < b.N; n++ {
		writers.WriteSchemaType(w, "Config", m, s, 0)
	}
}
