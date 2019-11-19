// Code generated by genfig plugin 'write_to_env' on 2019-11-19T10:45:54+01:00; DO NOT EDIT.

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	_ = os.Setenv
	_ = fmt.Sprintf
	_ = json.Marshal
)

func (c *Config) WriteToEnv() {
	var buf []byte
	_ = buf

	_ = os.Setenv("APIS_GOOGLE_URI", fmt.Sprintf("%v", c.Apis.Google.Uri))

	_ = os.Setenv("DB_PASS", fmt.Sprintf("%v", c.Db.Pass))

	_ = os.Setenv("DB_URI", fmt.Sprintf("%v", c.Db.Uri))

	_ = os.Setenv("DB_USER", fmt.Sprintf("%v", c.Db.User))

	_ = os.Setenv("LONGDESC_DE", fmt.Sprintf("%v", c.LongDesc.De))

	_ = os.Setenv("LONGDESC_EN", fmt.Sprintf("%v", c.LongDesc.En))

	_ = os.Setenv("PROJECT", fmt.Sprintf("%v", c.Project))

	_ = os.Setenv("RANDOMIZER_THRESHOLD", fmt.Sprintf("%v", c.Randomizer.Threshold))

	buf, _ = json.Marshal(c.Secrets)
	_ = os.Setenv("SECRETS", string(buf))

	_ = os.Setenv("SERVER_HOST", fmt.Sprintf("%v", c.Server.Host))

	_ = os.Setenv("SERVER_PORT", fmt.Sprintf("%v", c.Server.Port))

	_ = os.Setenv("VERSION", fmt.Sprintf("%v", c.Version))

	_ = os.Setenv("WIP", fmt.Sprintf("%v", c.Wip))

}
