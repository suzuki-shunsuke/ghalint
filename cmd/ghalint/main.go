package main

import (
	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

var version = ""

func main() {
	urfave.Main("ghalint", version, cli.Run)
}
