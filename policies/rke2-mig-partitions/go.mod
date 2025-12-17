module github.com/SUSE/openplatform-kubewarden-policies/policies/rke2-mig-partitions

go 1.22

toolchain go1.24.6

require (
	github.com/kubewarden/k8s-objects v1.32.0-kw1
	github.com/kubewarden/policy-sdk-go v0.12.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-openapi/strfmt v0.25.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/wapc/wapc-guest-tinygo v0.3.3 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-openapi/strfmt => github.com/kubewarden/strfmt v0.1.3
