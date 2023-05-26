package bpmetadata

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func buildUIInputFromVariables(vars []BlueprintVariable, input *BlueprintUIInput) {
	if input.Variables == nil {
		input.Variables = make(map[string]*DisplayVariable)
	}

	for _, v := range vars {
		_, hasDisplayVar := input.Variables[v.Name]
		if hasDisplayVar {
			continue
		}

		input.Variables[v.Name] = &DisplayVariable{
			Name:  v.Name,
			Title: createTitleFromName(v.Name),
		}
	}
}

func createTitleFromName(name string) string {
	nameSplit := strings.Split(name, "_")
	var titleSplit []string
	for _, n := range nameSplit {
		titleSplit = append(titleSplit, cases.Title(language.Und, cases.NoLower).String(n))
	}

	return strings.Join(titleSplit, " ")
}
