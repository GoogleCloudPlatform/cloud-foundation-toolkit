package validation

import (
	"context"
	"path/filepath"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	"github.com/pkg/errors"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func validateAssets(ctx context.Context, assets []cai.Asset, policyRootPath string) (*validator.AuditResponse, error) {
	return validateAssetsWithLibrary(ctx, assets,
		[]string{filepath.Join(policyRootPath, "policies")},
		filepath.Join(policyRootPath, "lib"))
}

func validateAssetsWithLibrary(ctx context.Context, assets []cai.Asset, policyPaths []string, policyLibraryDir string) (*validator.AuditResponse, error) {
	valid, err := gcv.NewValidator(policyPaths, policyLibraryDir)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcv validator")
	}

	auditResult := &validator.AuditResponse{}
	for i := range assets {
		asset := &validator.Asset{}
		if err := protoViaJSON(assets[i], asset); err != nil {
			return nil, errors.Wrapf(err, "converting asset %s to proto", assets[i].Name)
		}

		violations, err := valid.ReviewAsset(ctx, asset)
		if err != nil {
			return nil, errors.Wrapf(err, "reviewing asset %s", asset)
		}
		auditResult.Violations = append(auditResult.Violations, violations...)
	}
	return auditResult, nil
}
