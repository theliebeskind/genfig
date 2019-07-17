// Code generated by genfig on 2019-07-17T21:38:33+02:00; DO NOT EDIT.

package config

import (
	"os"
)

// Current is the current config, selected by the curren env and
// updated by the availalbe env vars
var Current *Config

// This init tries to retrieve the current environment via the
// common env var 'ENV' and applies activated plugins
func init() {
	Current, _ = Get(os.Getenv("ENV"))

	// calling plugin env_updater
	Current.UpdateFromEnv()

	// calling plugin substitutor
	Current.Substitute()

}
