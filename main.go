package main

import (
	"os"

	"github.com/mauromedda/vauth/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
