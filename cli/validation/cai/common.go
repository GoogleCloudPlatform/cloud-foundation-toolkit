package cai

// DM resource representation
type Resource struct {
	Project    string
	Name       string
	Type       string
	Properties map[string]interface{}
}
