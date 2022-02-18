module github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli

go 1.16

require (
	cloud.google.com/go/asset v1.0.1
	cloud.google.com/go/storage v1.18.2
	github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test v0.0.0-20211207170406-bea525b1cf52
	github.com/GoogleCloudPlatform/config-validator v0.0.0-20211122204404-f3fd77c5c355
	github.com/briandowns/spinner v1.16.0
	github.com/gammazero/workerpool v1.1.2
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jedib0t/go-pretty/v6 v6.2.4
	github.com/mitchellh/go-testing-interface v1.14.2-0.20210217184823-a52172cd2f64
	github.com/open-policy-agent/opa v0.37.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/open-policy-agent/gatekeeper v3.0.4-beta.2+incompatible => github.com/open-policy-agent/gatekeeper v0.0.0-20210409021048-9b5e4cfe5d7e
