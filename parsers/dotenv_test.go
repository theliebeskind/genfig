package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/theliebeskind/genfig/parsers"
)

const (
	complexDotenv = `
A=b
C_D=1
C_E=2
F=[2,"3","g"]
G=0.5
`
)

var (
	complexDotenvResult = map[string]interface{}{"a": "b", "c": map[string]interface{}{"d": int64(1), "e": int64(2)}, "f": []interface{}{float64(2), "3", "g"}, "g": 0.5}
)

func Test_Dotenv(t *testing.T) {
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
		{"double occurency map on basic", args{[]byte("A=1\nA_A=2")}, nil, true},
		{"double occurency basic on map", args{[]byte("A_A=2\nA=1")}, nil, true},
		{"nested double occurency", args{[]byte("A_A_A=2\nA_A=1")}, nil, true},
		{"valid dotenv", args{[]byte("A=1")}, map[string]interface{}{"a": int64(1)}, false},
		{"also valid dotenv", args{[]byte("A: 1")}, map[string]interface{}{"a": int64(1)}, false},
		{"complex dotenv", args{[]byte(complexDotenv)}, complexDotenvResult, false},
	}
	s := DotenvStrategy{}
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
