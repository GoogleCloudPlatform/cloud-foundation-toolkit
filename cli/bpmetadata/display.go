package bpmetadata

import "strings"

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
		titleSplit = append(titleSplit, strings.Title(n))
	}

	return strings.Join(titleSplit, " ")
}
