module github.com/SUSE/openplatform-kubewarden-policies/policies/pod-mig-partitions

go 1.25

require (
	github.com/kubewarden/k8s-objects v1.32.0-kw1
	github.com/kubewarden/policy-sdk-go v0.13.1
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/go-openapi/strfmt v0.25.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/wapc/wapc-guest-tinygo v0.3.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace github.com/go-openapi/strfmt => github.com/kubewarden/strfmt v0.1.3
