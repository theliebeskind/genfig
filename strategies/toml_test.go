package strategies_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/theliebeskind/genfig/strategies"
)

const (
	complexToml = `
	a = "b"
	f = [
		"2",
		"3",
		"g"
	]
	
	[c]
	d = 1
	e = 2
`
)

var (
	complexTomlResult = map[string]interface{}{"a": "b", "c": map[string]interface{}{"d": int64(1), "e": int64(2)}, "f": []interface{}{"2", "3", "g"}}
)

func Test_Toml(t *testing.T) {
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
		{"valid toml", args{[]byte("a=1")}, map[string]interface{}{"a": int64(1)}, false},
		{"complex toml", args{[]byte(complexToml)}, complexTomlResult, false},
	}
	s := TomlStrategy{}
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
