package bpmetadata

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
			Name: v.Name,
		}
	}
}
