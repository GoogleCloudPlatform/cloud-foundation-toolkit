package bpmetadata

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func buildUIInputFromVariables(vars []*BlueprintVariable, input *BlueprintUIInput) {
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

// mergeExistingAltDefaults merges existing alt_defaults from an old BlueprintUIInput into a new one,
// preserving manually authored alt_defaults.
func mergeExistingAltDefaults(newInput, existingInput *BlueprintUIInput) {
	if existingInput == nil {
		return // Nothing to merge if existingInput is nil
	}

	for i, variable := range newInput.Variables {
		for _, existingVariable := range existingInput.Variables {
			if variable.Name == existingVariable.Name && existingVariable.AltDefaults != nil {
				newInput.Variables[i].AltDefaults = existingVariable.AltDefaults
			}
		}
	}
}
