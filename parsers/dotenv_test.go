package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/theliebeskind/go-genfig/parsers"
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
		{"valid dotenv", args{[]byte("A=1")}, map[string]interface{}{"a": int64(1)}, false},
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
