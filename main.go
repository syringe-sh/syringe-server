package main

import (
	"github.com/nixpig/syringe.sh/internal/cli"
)

func main() {
	if err := cli.Cmd().Execute(); err != nil {
		// log.Fatal(err)
	}
}
