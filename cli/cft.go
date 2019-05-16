package main

import (
	"github.com/kopachevsky/cloud-foundation-toolkit/cli/cmd"
)

var Version string

func main() {
	cmd.Version = Version
	cmd.Execute()
}
