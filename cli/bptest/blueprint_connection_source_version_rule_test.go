package bptest

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"github.com/stretchr/testify/assert"
)

func TestBlueprintConnectionSourceVersionRule(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		expectErr    bool
		errorMessage string
	}{
		{
			name:      "Valid version - no equal sign expect version",
			version:   "1.2.3",
			expectErr: false,
		},
		{
			name:      "Valid version - expect version",
			version:   "=1.2.3",
			expectErr: false,
		},
		{
			name:      "Valid version - pessimistic constraint",
			version:   "~> 6.0",
			expectErr: false,
		},
		{
			name:      "Valid version - minimal version",
			version:   ">= 0.13.7",
			expectErr: false,
		},
		{
			name:      "Valid version - range interval",
			version:   ">= 0.13.7, < 2.0.0",
			expectErr: false,
		},
		{
			name:         "Invalid version - random string",
			version:      "invalid_version",
			expectErr:    true,
			errorMessage: "invalid_version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := &bpmetadata.BlueprintMetadata{
				Spec: &bpmetadata.BlueprintMetadataSpec{
					Interfaces: &bpmetadata.BlueprintInterface{
						Variables: []*bpmetadata.BlueprintVariable{
							{
								Connections: []*bpmetadata.BlueprintConnection{
									{
										Source: &bpmetadata.ConnectionSource{
											Source:  "example/source",
											Version: tt.version,
										},
									},
								},
							},
						},
					},
				},
			}

			ctx := lintContext{
				metadata: metadata,
			}

			rule := &BlueprintConnectionSourceVersionRule{}
			err := rule.check(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
