[![Kubewarden Policy Repository](https://github.com/kubewarden/community/blob/main/badges/kubewarden-policies.svg)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#policy-scope)
[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# harvester-restricted-network

## Settings

| Field                                                                                 | Description                       |
|---------------------------------------------------------------------------------------|-----------------------------------|
| namespaceVLANBindings <br> map[string, [NamespaceVLANBinding](#namespaceVLANBinding)] | A map of namespace VLAN bindings. |

### NamespaceVLANBinding

| Field                  | Description                            |
|------------------------|----------------------------------------|
| namespace <br/> string | The namespace.                         |
| vlan <br/> int         | The VLAN for the Harvester VM Network. |


## Specifications

1. All bound namespaces must use their respective bound VLANs.
2. All bound VLANs must use their respective bound namespaces.
3. Any namespace or VLAN that isn't bound, is unrestricted.

## Example

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: restricted-network-policy-1
spec:
  module: registry://ghcr.io/suse/openplatform-kubewarden-policies/harvester-restricted-network:latest
  rules:
    - apiGroups: ["kubevirt.io"]
      apiVersions: ["v1"]
      resources: ["virtualmachines"]
      operations: ["CREATE", "UPDATE"]
  settings:
    namespaceVLANBindings:
      - namespace:  test-restricted-1
        vlan:       42
      - namespace:  test-restricted-2
        vlan:       1337
  mutating: false
  policyServer: default
```

Here would be the result of the above policy.

| namespace         | network | Result |
|-------------------|---------|--------|
| test-restricted-1 | 42      | ALLOW  |
| test-restricted-2 | 1337    | ALLOW  |
| random-namespace  | 100     | ALLOW  |
| test-restricted-1 | 1337    | REJECT |
| test-restricted-2 | 42      | REJECT |
| random-namespace  | 42      | REJECT |
| random-namespace  | 1337    | REJECT |
| test-restricted-1 | 100     | REJECT |
| test-restricted-2 | 100     | REJECT |
