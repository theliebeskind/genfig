//go:generate go run ../cmd ../fixtures/default.yml ../fixtures/*.yaml ../fixtures/*.toml ../fixtures/.env*

package main

import (
	"fmt"

	"github.com/theliebeskind/genfig/example/config"
)

func main() {
	fmt.Println(config.Current.Version)
	fmt.Println(config.Current.Randomizer.Threshold)
	fmt.Println(config.Current.Secrets)
	fmt.Println(config.Current.Server.Port)
	fmt.Println(config.Current.Db)
}
