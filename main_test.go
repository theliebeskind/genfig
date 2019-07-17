package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	fixturesDir, _ = filepath.Abs("./fixtures")
	origArgs       = os.Args
)

func Test_run(t *testing.T) {
	dir, _ := ioutil.TempDir("", "genfig")
	out, _ := ioutil.TempDir("", "genfig")
	defer os.RemoveAll(dir)
	defer os.RemoveAll(out)

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(dir)
	require.NoError(t, err)
	fmt.Printf("CWD is now %s\n", dir)

	configsDir := filepath.Join(fixturesDir, "configs")

	tests := []struct {
		name        string
		args        []string
		shouldPanic bool
	}{
		{"no args", []string{}, true},
		{"help", []string{"--help"}, false},
		{"without dir, no config files", []string{"*"}, true},
		{"with dir, no config files", []string{"-dir", out, "*"}, true},
		{"without dir, valid config files", []string{configsDir + "/default.yml", configsDir + "/development.yml"}, false},
		{"with dir, valid config files", []string{"-dir", out, configsDir + "/default.yml", configsDir + "/development.yml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			args := append(origArgs[:1], tt.args...)
			os.Args = args
			if tt.shouldPanic {
				require.Panics(t, run)
			} else {
				require.NotPanics(t, run)
			}
			os.Args = origArgs
			os.RemoveAll(dir)
		})
	}
}
