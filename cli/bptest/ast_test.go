package bptest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTestFuncsFromFile(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   []string
		errMsg string
	}{
		{
			name: "simple",
			data: `package test

import "testing"

func TestA(t *testing.T) {
}
`,
			want: []string{"TestA"},
		},
		{
			name: "multiple",
			data: `package test

import "testing"

const ShouldNotErr = "foo"

func TestA(t *testing.T) {
}

func TestB(t *testing.T) {
}

func OtherHelper(t *testing.T) {
}
`,
			want: []string{"TestA", "TestB"},
		},
		{
			name: "empty",
			data: `package test
`,
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			filePath, cleanup := writeTmpFile(t, tt.data)
			defer cleanup()
			got, err := getTestFuncsFromFile(filePath)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.ElementsMatch(tt.want, got)
			}
		})
	}
}

func writeTmpFile(t *testing.T, data string) (string, func()) {
	assert := assert.New(t)
	f, err := os.CreateTemp("", "*.go")
	assert.NoError(err)
	cleanup := func() { os.Remove(f.Name()) }
	_, err = f.Write([]byte(data))
	assert.NoError(err)
	f.Close()
	return f.Name(), cleanup
}
