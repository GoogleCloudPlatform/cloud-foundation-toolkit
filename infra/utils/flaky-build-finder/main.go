package main

import (
	"flag"
	"log"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/utils/fbf/cmd"
)

func main() {
	startTime := flag.String("start-time", "", "Time to start computing flakes in form MM-DD-YYYY")
	endTime := flag.String("end-time", "", "Time to stop computing flakes in form MM-DD-YYYY")
	projectID := flag.String("project-id", "", "Project ID")
	verbose := flag.Bool("verbose", false, "Display detailed table with flaky build IDs")
	flag.Parse()

	ftf, err := cmd.NewFlakyFinder(*startTime, *endTime, *projectID, *verbose)
	if err != nil {
		log.Fatalf("error initializing flaky finder: %v", err)
	}
	err = ftf.ComputeFlakes()
	if err != nil {
		log.Fatalf("error computing flakes: %v", err)
	}
	ftf.Render()
}
