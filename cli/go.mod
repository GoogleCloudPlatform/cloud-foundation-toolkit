module github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli

go 1.16

require (
	cloud.google.com/go/asset v1.10.1
	cloud.google.com/go/storage v1.28.1
	github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test v0.4.0
	github.com/GoogleCloudPlatform/config-validator 1d72524ea1b8
	github.com/briandowns/spinner v1.20.0
	github.com/fatih/color v1.13.0
	github.com/gammazero/workerpool v1.1.3
	github.com/go-git/go-git/v5 v5.5.1
	github.com/golang/protobuf v1.5.2
	github.com/gomarkdown/markdown 663e2500819c
	github.com/google/go-cmp v0.5.9
	github.com/google/go-github/v48 v48.2.0
	github.com/hashicorp/hcl/v2 v2.15.0
	github.com/hashicorp/terraform-config-inspect 81db043ad408
	github.com/iancoleman/strcase v0.2.0
	github.com/inconshreveable/log15 555555054819
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jedib0t/go-pretty/v6 v6.4.3
	github.com/manifoldco/promptui v0.9.0
	github.com/migueleliasweb/go-github-mock v0.0.13
	github.com/mitchellh/go-testing-interface 2d9075ca8770
	github.com/open-policy-agent/gatekeeper v3.0.4-beta.2+incompatible // indirect
	github.com/open-policy-agent/opa v0.34.2
	github.com/otiai10/copy v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.14.0
	github.com/stretchr/testify v1.8.1
	github.com/zclconf/go-cty v1.10.0 // indirect
	golang.org/x/oauth2 v0.3.0
	google.golang.org/api v0.105.0
	google.golang.org/genproto 3c3c17ce83e6
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/kustomize/kyaml v0.13.10
	sigs.k8s.io/yaml v1.3.0
)

replace github.com/open-policy-agent/gatekeeper v3.0.4-beta.2+incompatible => github.com/open-policy-agent/gatekeeper v0.0.0-20210409021048-9b5e4cfe5d7e
