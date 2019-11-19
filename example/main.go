//go:generate rm -rf config
//go:generate go run ../ --dir config ../fixtures/configs/default.yml ../fixtures/configs/*.yaml ../fixtures/configs/*.json ../fixtures/configs/*.toml ../fixtures/configs/.env*

package main

import (
	"fmt"

	"github.com/thlcodes/genfig/example/config"
)

func main() {
	fmt.Println(config.Current.Version)
	fmt.Println(config.Current.Randomizer.Threshold)
	fmt.Println(config.Current.Secrets)
	fmt.Println(config.Current.Server.Port)
	fmt.Println(config.Current.Db)

	fmt.Println(config.Envs.Test.Version)

	fmt.Println(config.Current.Map())
}
