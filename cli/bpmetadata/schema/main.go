package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	log "github.com/inconshreveable/log15"
)

// bpmetadata log15 handler
var Log = log.New()

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"USAGE: %s [-output=PATH]\n",
			path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	output := flag.String("output", "", "output path for generating the JSON schema definition")
	flag.Parse()

	os.Exit(process(*output))
}

func process(output string) int {
	// get the working directory for the command
	wdPath, err := os.Getwd()
	if err != nil {
		Log.Error("error getting working dir", "err", err)
		return 1
	}

	if err := generateSchemaFile(output, wdPath); err != nil {
		Log.Error("error generating schema", "err", err)
		return 1
	}

	return 0
}
