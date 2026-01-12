[![Kubewarden Policy Repository](https://github.com/kubewarden/community/blob/main/badges/kubewarden-policies.svg)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#policy-scope)
[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# harvester-pci-devices

This policy guards against VMs attaching PCI Devices (e.g., GPUs) without permission.

## Settings

| Field                                                                                      | Description                             |
|--------------------------------------------------------------------------------------------|-----------------------------------------|
| namespaceDeviceindings <br> map[string, [NamespaceDeviceBinding](#namespaceDeviceBinding)] | A map of Harvester PCI Device bindings. |

### NamespaceDeviceBinding

| Field                  | Description               |
|------------------------|---------------------------|
| namespace <br/> string | The namespace.            |
| device <br/> string    | The ID of the PCI device. |


## Specifications

1. You should be able to create a VM without a PCI Device
2. You should not be able to bind a VM to a PCI Device not allocated to its namespace.

## Example

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: harvester-pci-policy-1
spec:
  module: harbor.op-prg2-0-dev-ingress.op.suse.org/op-portal/kubewarden-policy:20
  rules:
    - apiGroups: ["kubevirt.io"]
      apiVersions: ["v1"]
      resources: ["virtualmachines"]
      operations: ["CREATE", "UPDATE"]
  settings:
    namespaceDeviceBindings:
      - namespace: test-ns-1
        device:  tekton27a-000001010
      - namespace: test-ns-2
        device:  tekton27b-000001010
  mutating: false  # or true if your policy mutates resources
  policyServer: default
```

Here would be the result of the above policy.

| Namespace        | PCI Device ID       | Result |
|------------------|---------------------|--------|
| test-ns-1        | tekton27a-000001010 | ALLOW  |
| test-ns-2        | tekton27b-000001010 | ALLOW  |
| test-ns-1        | tekton27b-000001010 | REJECT |
| test-ns-2        | tekton27a-000001010 | REJECT |
| random-namespace | tekton27a-000001010 | REJECT |
| random-namespace | tekton27b-000001010 | REJECT |
