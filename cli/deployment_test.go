package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	path := filepath.Join("testdata", "config", "simple.yaml")
	config := NewConfig(path)
	if config == nil {
		t.Errorf("Config is nil")
	}
}

func TestFindAllOutRefs(t *testing.T) {
	executeFindAllOutRefsAndAssert(
		`network: $(out.project1.deployment1.resource1.output_name1)
                    other: $(out.resource2.output2)`,
		[]string{"$(out.project1.resource1.output1)", "$(out.resource2.output2)"},
		t)
	// empty string
	executeFindAllOutRefsAndAssert("", nil, t)
	// invalid notation
	executeFindAllOutRefsAndAssert("${out1.account.project.resource.output", nil, t)
}

func executeFindAllOutRefsAndAssert(yaml_string string, expected []string, t *testing.T) {
	config := &Config{yaml_string: yaml_string}
	actual := config.findAllOutRefs()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %v, expected: %v.", actual, expected)
	}
}

func TestNewConfigGraph(t *testing.T) {

}
