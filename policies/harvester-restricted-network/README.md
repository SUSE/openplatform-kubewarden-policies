# harvester-restricted-netwok

This policy guards against harvester creating a network for a restricted VLAN in unauthorized namespaces.

**Example policy:**

```
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: restricted-vlan-policy-1
spec:
  module: harvester-restricted-network:0.1.0
  rules:
    - apiGroups: ["k8s.cni.cncf.io"]
      apiVersions: ["v1"]
      resources: ["network-attachment-definitions"]
      operations: ["CREATE", "UPDATE"]
  settings:
    namespaceVLANBindings:
      - namespace:  test-restricted-1-network-1
        vlan:       42
      - namespace:  test-restricted-2-network-1
        vlan:       1337
  mutating: false  # or true if your policy mutates resources
  policyServer: default
```

**Specifications:**

1. All bound namespaces must use their respective bound VLANs.
2. All bound VLANs must use their respective bound namespaces.
3. Any namespace or VLAN that isn't bound, is unrestricted.

**Examples:**

The following examples are with the example policy above, with a random non-restricted VLAN being 100.

| namespace                   | network | Result |
|-----------------------------|---------|--------|
| test-restricted-1-network-1 | 42      | ALLOW  |
| test-restricted-2-network-1 | 1337    | ALLOW  |
| random-namespace            | 100     | ALLOW  |
| test-restricted-1-network-1 | 1337    | REJECT |
| test-restricted-2-network-2 | 42      | REJECT |
| random-namespace            | 42      | REJECT |
| random-namespace            | 1337    | REJECT |
| test-restricted-1-network-1 | 100     | REJECT |
| test-restricted-2-network-2 | 100     | REJECT |
