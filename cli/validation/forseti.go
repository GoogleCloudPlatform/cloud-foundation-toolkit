package validation

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation/cai"
)

func validateAssets(assets []cai.Asset, policyRootPath string) (*validator.AuditResponse, error) {
	return validateAssetsWithLibrary(assets,
		filepath.Join(policyRootPath, "policies"),
		filepath.Join(policyRootPath, "lib"))
}

func validateAssetsWithLibrary(assets []cai.Asset, policyPath, policyLibraryDir string) (*validator.AuditResponse, error) {
	valid, err := gcv.NewValidator(
		gcv.PolicyPath(policyPath),
		gcv.PolicyLibraryDir(policyLibraryDir),
	)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcv validator")
	}

	pbAssets := make([]*validator.Asset, len(assets))
	for i := range assets {
		pbAssets[i] = &validator.Asset{}
		if err := protoViaJSON(assets[i], pbAssets[i]); err != nil {
			return nil, errors.Wrapf(err, "converting asset %s to proto", assets[i].Name)
		}
	}

	if err := valid.AddData(&validator.AddDataRequest{
		Assets: pbAssets,
	}); err != nil {
		return nil, errors.Wrap(err, "adding data to validator")
	}

	auditResult, err := valid.Audit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "auditing")
	}

	return auditResult, nil
}
