package bptest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAndGetStage(t *testing.T) {
	tests := []struct {
		name   string
		stage  string
		want   string
		errMsg string
	}{
		{
			name:  "valid name",
			stage: "init",
			want:  "init",
		},
		{
			name:  "alias name",
			stage: "create",
			want:  "init",
		},
		{
			name:  "valid name no alias",
			stage: "verify",
			want:  "verify",
		},
		{
			name:   "invalid name",
			stage:  "foo",
			errMsg: "invalid stage name foo - one of [\"init\" \"apply\" \"verify\" \"teardown\"] expected",
		},
		{
			name:  "empty (all stages)",
			stage: "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got, err := validateAndGetStage(tt.stage)
			if tt.errMsg != "" {
				assert.NotNil(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}
