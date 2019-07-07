package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fixturesDir = "../fixtures/"
)

var (
	origArgs = os.Args
)

func Test_exec(t *testing.T) {
	dir := "out"

	tests := []struct {
		name        string
		args        []string
		shouldPanic bool
	}{
		{"no args", []string{}, true},
		{"without dir, no config files", []string{"*"}, true},
		{"with dir, no config files", []string{"-dir", dir, "*"}, true},
		{"without dir, valid config files", []string{fixturesDir + "default.yml", fixturesDir + "development.yml"}, false},
		{"with dir, valid config files", []string{"-dir", dir, fixturesDir + "default.yml", fixturesDir + "development.yml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append(origArgs[:1], tt.args...)
			os.Args = args
			if tt.shouldPanic {
				require.Panics(t, exec)
			} else {
				require.NotPanics(t, exec)
			}
			os.Args = origArgs
			os.RemoveAll(dir)
		})
	}
}
