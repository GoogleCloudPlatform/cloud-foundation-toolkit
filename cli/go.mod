module github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli

go 1.16

require (
	cloud.google.com/go/asset v1.11.1
	cloud.google.com/go/storage v1.29.0
	github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test v0.5.0
	github.com/GoogleCloudPlatform/config-validator 0dfa3040e8a4
	github.com/briandowns/spinner v1.22.0
	github.com/fatih/color v1.14.1
	github.com/gammazero/workerpool v1.1.3
	github.com/go-git/go-git/v5 v5.6.0
	github.com/golang/protobuf v1.5.3
	github.com/gomarkdown/markdown 3238e54d4819
	github.com/google/go-cmp v0.5.9
	github.com/google/go-github/v50 v50.1.0
	github.com/hashicorp/hcl/v2 v2.16.1
	github.com/hashicorp/terraform-config-inspect d7dec65d5f3a
	github.com/iancoleman/strcase v0.2.0
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jedib0t/go-pretty/v6 v6.4.6
	github.com/manifoldco/promptui v0.9.0
	github.com/migueleliasweb/go-github-mock v0.0.16
	github.com/mitchellh/go-testing-interface 2d9075ca8770
	github.com/open-policy-agent/gatekeeper v3.0.4-beta.2+incompatible // indirect
	github.com/open-policy-agent/opa v0.34.2
	github.com/otiai10/copy v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.2
	github.com/zclconf/go-cty v1.10.0 // indirect
	golang.org/x/oauth2 v0.5.0
	google.golang.org/api v0.111.0
	google.golang.org/genproto 7f2fa6fef1f4
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/kustomize/kyaml v0.14.0
	sigs.k8s.io/yaml v1.3.0
)

replace github.com/open-policy-agent/gatekeeper v3.0.4-beta.2+incompatible => github.com/open-policy-agent/gatekeeper v0.0.0-20210409021048-9b5e4cfe5d7e
