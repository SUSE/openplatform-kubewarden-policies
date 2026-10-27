module github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network-vm

go 1.25

require (
	github.com/francoispqt/onelog v0.0.0-20190306043706-8c2bb31b10a4
	github.com/kubewarden/policy-sdk-go v0.13.1
	github.com/stretchr/testify v1.12.1
	github.com/wapc/wapc-guest-tinygo v0.3.3
)

require (
	github.com/francoispqt/gojay v0.0.0-20181220093123-f2cc13a668ca // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/kubewarden/k8s-objects v1.29.0-kw1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace github.com/go-openapi/strfmt => github.com/kubewarden/strfmt v0.1.3
