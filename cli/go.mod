module github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli

go 1.16

require (
	cloud.google.com/go/asset v1.0.1
	cloud.google.com/go/storage v1.18.2
	github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test v0.0.0-20220307174651-21d0dee0c8ea
	github.com/GoogleCloudPlatform/config-validator v0.0.0-20211122204404-f3fd77c5c355
	github.com/briandowns/spinner v1.16.0
	github.com/fatih/color v1.13.0
	github.com/gammazero/workerpool v1.1.2
	github.com/go-git/go-git/v5 v5.4.2
	github.com/golang/protobuf v1.5.2
	github.com/gomarkdown/markdown v0.0.0-20220905174103-7b278df48cfb
	github.com/google/go-cmp v0.5.9
	github.com/google/go-github/v47 v47.1.0
	github.com/hashicorp/hcl/v2 v2.14.0
	github.com/hashicorp/terraform-config-inspect v0.0.0-20211115214459-90acf1ca460f
	github.com/iancoleman/strcase v0.2.0
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jedib0t/go-pretty/v6 v6.2.4
	github.com/manifoldco/promptui v0.9.0
	github.com/migueleliasweb/go-github-mock v0.0.12
	github.com/mitchellh/go-testing-interface v1.14.2-0.20210217184823-a52172cd2f64
	github.com/open-policy-agent/gatekeeper v3.0.4-beta.2+incompatible // indirect
	github.com/open-policy-agent/opa v0.34.2
	github.com/otiai10/copy v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.8.0
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	google.golang.org/api v0.58.0
	google.golang.org/genproto v0.0.0-20211129164237-f09f9a12af12
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	sigs.k8s.io/kustomize/kyaml v0.13.9
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/open-policy-agent/gatekeeper v3.0.4-beta.2+incompatible => github.com/open-policy-agent/gatekeeper v0.0.0-20210409021048-9b5e4cfe5d7e
