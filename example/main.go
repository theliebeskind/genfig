//go:generate go run ../cmd ../fixtures/default.yml ../fixtures/*.yaml ../fixtures/*.toml ../fixtures/.env*

package main

import (
	"fmt"

	"github.com/theliebeskind/genfig/example/config"
)

func main() {
	fmt.Println(config.Default.Version)
}
