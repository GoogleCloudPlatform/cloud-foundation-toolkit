package bptest

import "fmt"

var stages = []string{"init", "plan", "apply", "verify", "teardown"}

var stagesWithAlias = map[string][]string{
	stages[0]: {"create"},
	stages[1]: {},
	stages[2]: {"converge"},
	stages[3]: {},
	stages[4]: {"destroy"},
}

// validateAndGetStage validates given stage and resolves to stage name if an alias is provided
func validateAndGetStage(s string) (string, error) {
	// empty stage is a special case for running all stages
	if s == "" {
		return "", nil
	}
	for stageName, aliases := range stagesWithAlias {
		if stageName == s {
			return stageName, nil
		}
		for _, alias := range aliases {
			if alias == s {
				return stageName, nil
			}
		}
	}
	return "", fmt.Errorf("invalid stage name %s - one of %+q expected", s, stages)
}
