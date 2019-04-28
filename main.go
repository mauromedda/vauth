package main

import (
	"os"

	"github.com/mauromedda/vauth/command"
)

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
