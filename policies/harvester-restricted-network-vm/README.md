[![Kubewarden Policy Repository](https://github.com/kubewarden/community/blob/main/badges/kubewarden-policies.svg)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#policy-scope)
[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# harvester-restricted-network-vm

This policy protects Harvester VM networks, by specifying which namespaces are allowed.

## Settings

| Field                                                                                          | Description                          |
|------------------------------------------------------------------------------------------------|--------------------------------------|
| namespaceNetworkBindings <br> map[string, [NamespaceNetworkBinding](#namespaceNetworkBinding)] | A map of namespace network bindings. |

### NamespaceNetworkBinding

| Field                  | Description                                                      |
|------------------------|------------------------------------------------------------------|
| namespace <br/> string | The namespace.                                                   |
| network <br/> string   | The Harvester VM Network in the format `namespace/network-name`. |

## Specifications

1. You should be able to create a VM with any of the specified combinations of namespace and network.
2. You should not be able to create a VM from any namespace or network that is in the settings, but the exact combination is not in the settings.
3. Any namespace or network that is not on the settings is not restricted

## Example

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: restricted-network-vm-policy-1
spec:
  module: registry://ghcr.io/suse/openplatform-kubewarden-policies/harvester-restricted-network-vm:latest
  rules:
    - apiGroups: ["kubevirt.io"]
      apiVersions: ["v1"]
      resources: ["virtualmachines"]
      operations: ["CREATE", "UPDATE"]
  settings:
    namespaceNetworkBindings:
      - namespace: test-restricted-1-network-1
        network:  test-restricted-1-network-1/network-1
      - namespace: test-restricted-2-network-1
        network:  test-restricted-1-network-1/network-1
      - namespace: test-restricted-3-network-3
        network:  test-restricted-3-network-3/network-3
  mutating: false
  policyServer: default
```

Here would be the result of the above policy.

| namespace                   | network                               | Result |
|-----------------------------|---------------------------------------|--------|
| test-restricted-1-network-1 | test-restricted-1-network-1/network-1 | ALLOW  |
| test-restricted-2-network-1 | test-restricted-1-network-1/network-1 | ALLOW  |
| test-restricted-3-network-3 | test-restricted-3-network-3/network-3 | ALLOW  |
| random-namespace            | random-network                        | ALLOW  |
| test-restricted-3-network-3 | test-restricted-1-network-1/network-1 | REJECT |
| test-restricted-1-network-1 | test-restricted-3-network-3/network-3 | REJECT |
| random-namespace            | test-restricted-1-network-1/network-1 | REJECT |
| random-namespace            | test-restricted-3-network-3/network-3 | REJECT |
| test-restricted-1-network-1 | random-network                        | REJECT |
| test-restricted-2-network-2 | random-network                        | REJECT |
| test-restricted-3-network-3 | random-network                        | REJECT |
