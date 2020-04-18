package main

import (
	"cloud-foundation-toolkit/config-connector/tests/ccs-test/cmd"
	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()

	cmd.Execute()
}
