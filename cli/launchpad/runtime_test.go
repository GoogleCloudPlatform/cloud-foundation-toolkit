package launchpad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyResource struct {
	headerYAML
	id string
}

func (d *dummyResource) resId() string                                  { return "Folder." + d.id }
func (d *dummyResource) validate() error                                { return nil }
func (d *dummyResource) kind() crdKind                                  { return crdKind(0) }
func (d *dummyResource) addToOrg(ao *assembledOrg) error                { return nil }
func (d *dummyResource) resolveReferences(refs []resourceHandler) error { return nil }

type registerResourceArg struct {
	src resourceHandler
	dst *referenceYAML
}

func TestAssembledOrg_registerResourceAddResourceRegistry(t *testing.T) {
	d1, d2 := &dummyResource{id: "1"}, &dummyResource{id: "2"}
	expAo1 := newAssembledOrg()
	expAo1.resourceMap[d1.resId()] = &resource{yaml: d1}

	expAo2 := newAssembledOrg()
	expAo2.resourceMap[d1.resId()] = &resource{yaml: d1}
	expAo2.resourceMap[d2.resId()] = &resource{yaml: d2}

	var testCases = []struct {
		name           string
		inputResources []resourceHandler
		expectedOutput *assembledOrg
	}{{
		"add_one_resource",
		[]resourceHandler{d1},
		expAo1,
	}, {
		"add_two_resources",
		[]resourceHandler{d1, d2},
		expAo2,
	}, {
		"add_same_resource_twice",
		[]resourceHandler{d1, d2, d1},
		expAo2,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := newAssembledOrg()
			for _, r := range tc.inputResources {
				assert.Nil(t, ao.registerResource(r, nil))
			}
			assert.Equal(t, tc.expectedOutput, ao, "assembleOrg is expected to be equal")
		})
	}
}

func TestAssembledOrg_registerResourceAddReferenceTracker(t *testing.T) {
	d1, d2, d3 := &dummyResource{id: "1"}, &dummyResource{id: "2"}, &dummyResource{id: "3"}
	expAo0 := newAssembledOrg() // Case 1 (no ref)
	expAo0.resourceMap[d1.resId()] = &resource{yaml: d1}

	expAo1 := newAssembledOrg() // Case 1 <- 2
	expAo1.resourceMap[d1.resId()] = &resource{yaml: d1, inRefs: []resourceHandler{d2}}
	expAo1.resourceMap[d2.resId()] = &resource{yaml: d2}

	expAo2 := newAssembledOrg() // Case 1 <- 2, 3
	expAo2.resourceMap[d1.resId()] = &resource{yaml: d1, inRefs: []resourceHandler{d2, d3}}
	expAo2.resourceMap[d2.resId()] = &resource{yaml: d2}
	expAo2.resourceMap[d3.resId()] = &resource{yaml: d3}

	expAoUndefined := newAssembledOrg() // Case 1 -> 000 (not exist)
	expAoUndefined.resourceMap[d1.resId()] = &resource{yaml: d1}
	expAoUndefined.resourceMap["Folder.000"] = &resource{inRefs: []resourceHandler{d1}}

	var testCases = []struct {
		name           string
		inputs         []registerResourceArg
		expectedOutput *assembledOrg
	}{{
		"no_reference",
		[]registerResourceArg{
			{d1, nil},
		},
		expAo0,
	}, {
		"single_reference",
		[]registerResourceArg{
			{d1, nil},
			{d2, &referenceYAML{"folder", "1"}},
		},
		expAo1,
	}, {
		"double_references",
		[]registerResourceArg{
			{d1, nil},
			{d2, &referenceYAML{"folder", "1"}},
			{d3, &referenceYAML{"folder", "1"}},
		},
		expAo2,
	}, {
		"duplicated_references",
		[]registerResourceArg{
			{d1, nil},
			{d2, &referenceYAML{"folder", "1"}},
			{d3, &referenceYAML{"folder", "1"}},
		},
		expAo2,
	}, {
		"undefined_reference",
		[]registerResourceArg{
			{d1, &referenceYAML{"folder", "000"}},
		},
		expAoUndefined,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := newAssembledOrg()
			for _, r := range tc.inputs {
				assert.Nil(t, ao.registerResource(r.src, r.dst))
			}
			assert.EqualValues(t, tc.expectedOutput.resourceMap, ao.resourceMap, "assembledOrg is expected to be equal")
		})
	}
}

func TestAssembledOrg_registerResourceInitializeOrg(t *testing.T) {
	d1, d2 := &dummyResource{id: "1"}, &dummyResource{id: "2"}
	expAo0 := newAssembledOrg() // no org init
	expAo1 := newAssembledOrg() // org init
	expAo1.org.Spec.Id = "1234"

	var testCases = []struct {
		name           string
		inputs         []registerResourceArg
		expectedOutput *assembledOrg
	}{{
		"no_org_init",
		[]registerResourceArg{{d1, nil}},
		expAo0,
	}, {
		"org_init",
		[]registerResourceArg{{d1, &referenceYAML{Organization.String(), "1234"}}},
		expAo1,
	}, {
		"org_init_multiple_times",
		[]registerResourceArg{
			{d1, &referenceYAML{Organization.String(), "1234"}},
			{d2, &referenceYAML{Organization.String(), "1234"}},
		},
		expAo1,
	}, {
		"org_init_conflict",
		[]registerResourceArg{
			{d1, &referenceYAML{Organization.String(), "1234"}},
			{d2, &referenceYAML{Organization.String(), "3333"}},
		},
		nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := newAssembledOrg()
			var err error
			for _, r := range tc.inputs {
				err = ao.registerResource(r.src, r.dst)
				if err != nil {
					break
				}
			}
			if err != nil && tc.expectedOutput == nil {
				assert.Equal(t, errConflictDefinition, err, "expect to have org id conflict error")
				return
			}
			assert.Equal(t, tc.expectedOutput.org, ao.org, "org is expected to be initialized the same")

			if len(tc.expectedOutput.resourceMap) == 0 {
				return
			}
			_, found := ao.resourceMap[Organization.String()+".1234"]
			assert.True(t, found, "org is expected to be registered after initialization")
		})
	}
}

func TestAssembledOrg_resolveReferencesNotFound(t *testing.T) {
	d2 := &dummyResource{id: "2"}
	ao := newAssembledOrg()
	_ = ao.registerResource(d2, &referenceYAML{"Folder", "1"}) // reference Folder.1

	err := ao.resolveReferences()
	assert.Equal(t, errUndefinedReference, err)
}
