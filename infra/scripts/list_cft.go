//usr/bin/env go run "$0" "$@"; exit "$?"

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v25/github"
	"github.com/kr/pretty"
)

func main() {

	client := github.NewClient(nil)

	orgs, _, _ := client.Search.Repositories(context.Background(), "willnorris", nil)

	pretty.Printf(orgs)
	fmt.Println("Hello")
	os.Exit(42)
}
