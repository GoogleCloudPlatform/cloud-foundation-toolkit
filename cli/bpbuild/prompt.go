package bpbuild

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

// promptSelect prompts a user to select a value from given items.
func promptSelect(label string, items []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Searcher: func(input string, index int) bool {
			return strings.Contains(items[index], input)
		},
		StartInSearchMode: true,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Selected build step: %s\n", result)
	return result
}
