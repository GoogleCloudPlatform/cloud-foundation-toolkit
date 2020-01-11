package launchpad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawingIntegration(t *testing.T) {
	var testCases = []struct {
		name       string
		inputYAMLs []string
		output     string
	}{{
		"basic_org",
		[]string{"../testdata/launchpad/folder/org_1.yaml"},
		`elements {
  group organization {
    card generic as org {
      name "organization"
    }
    card generic as group1 {
      name "group1"
    }
    card generic as group2 {
      name "group2"
    }
  }
}
paths {
  org-->group1
  org-->group2
}
`,
	}, {
		"nested_org",
		[]string{"../testdata/launchpad/folder/org_3_nested.yaml"},
		`elements {
  group organization {
    card generic as org {
      name "organization"
    }
    card generic as group1 {
      name "group1"
    }
    card generic as group1_dev {
      name "Dev"
    }
    card generic as group1_uat {
      name "UAT"
    }
    card generic as group1_prod {
      name "Production"
    }
    card generic as group2 {
      name "group2"
    }
    card generic as group2_dev {
      name "Dev"
    }
    card generic as group2_uat {
      name "UAT"
    }
    card generic as group2_prod {
      name "Production"
    }
  }
}
paths {
  org-->group1
  group1-->group1_dev
  group1-->group1_uat
  group1-->group1_prod
  org-->group2
  group2-->group2_dev
  group2-->group2_uat
  group2-->group2_prod
}
`,
	}, {
		"nested_deep",
		[]string{"../testdata/launchpad/folder/org_4_deep.yaml"},
		`elements {
  group organization {
    card generic as org {
      name "organization"
    }
    card generic as group1 {
      name "group1"
    }
    card generic as group1_nest1 {
      name "Nest 1.1"
    }
    card generic as group1_nest1_nest2 {
      name "Nest 1.1.2"
    }
    card generic as group2 {
      name "group2"
    }
  }
}
paths {
  org-->group1
  group1-->group1_nest1
  group1_nest1-->group1_nest1_nest2
  org-->group2
}
`,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resources := loadResources(tc.inputYAMLs)
			assembled := assembleResourcesToOrg(resources)

			drawing, err := assembled.makeDiagram()
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tc.output, drawing.String())
		})
	}
}
