package main

import (
	"github.com/fynntang/staking-deposit/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
