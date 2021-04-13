package main

import (
	"flag"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/infra/developer-tools/build/scripts/module-swapper/cmd"
)

func main() {
	workDir := flag.String("workdir", "", "Absolute path to root module where examples should be swapped. Defaults to working directory")
	subModulesDir := flag.String("submods-path", "modules", "Path to a submodules if any that maybe referenced. Defaults to working dir/modules")
	examplesDir := flag.String("examples-path", "examples", "Path to examples that should be swapped. Defaults to cwd/examples")
	moduleRegistrySuffix := flag.String("registry-suffix", "google", "Module registry suffix")
	restore := flag.Bool("restore", false, "Restores disabled modules")
	flag.Parse()
	rootPath := *workDir
	// if no workDir specified default to current working directory
	if rootPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Unable to get cwd: %v", err)
		}
		rootPath = cwd
	}
	cmd.SwapModules(rootPath, *moduleRegistrySuffix, *subModulesDir, *examplesDir, *restore)
}
