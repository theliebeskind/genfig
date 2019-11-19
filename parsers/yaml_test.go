package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/thclodes/genfig/parsers"
)

const (
	complexYaml = `
a: b
c:
  d: 1
  e: 2
f:
  - 2
  - "3"
  - g	
`
	complexJson = `
{
	"a": "b",
	"c": {
		"d": 1,
		"e": 2
	},
	"f": [
		2,  
		"3",
		"g"
	]
}
`
)

var (
	complexYamlResult = map[string]interface{}{"a": "b", "c": map[interface{}]interface{}{"d": 1, "e": 2}, "f": []interface{}{2, "3", "g"}}
)

func Test_Yaml(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{"empty data", args{}, nil, true},
		{"invalid data", args{[]byte("foobar´?")}, nil, true},
		{"valid yaml", args{[]byte("a: 1")}, map[string]interface{}{"a": 1}, false},
		{"vaild json", args{[]byte(`{"a": 1}`)}, map[string]interface{}{"a": 1}, false},
		{"complex yaml", args{[]byte(complexYaml)}, complexYamlResult, false},
		{"complex json", args{[]byte(complexJson)}, complexYamlResult, false},
	}
	s := YamlStrategy{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Parse(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
