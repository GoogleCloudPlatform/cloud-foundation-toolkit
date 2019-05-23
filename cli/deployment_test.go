package main

import (
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	path := filepath.Join("testdata", "config", "simple.yaml")
	config := NewConfig(path)
	if config == nil {
		t.Errorf("Config is nil")
	}
}
