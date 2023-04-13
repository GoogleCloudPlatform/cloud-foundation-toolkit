package bpmetadata

func buildUIInputFromVariables(vars []BlueprintVariable, input *BlueprintUIInput) {
	if input.DisplayVariables == nil {
		input.DisplayVariables = make(map[string]*DisplayVariable)
	}

	for _, v := range vars {
		_, hasDisplayVar := input.DisplayVariables[v.Name]
		if hasDisplayVar {
			continue
		}

		input.DisplayVariables[v.Name] = &DisplayVariable{
			Name: v.Name,
		}
	}
}
