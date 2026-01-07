# harvester-restricted-network-vm

This policy guards against VMs being deployed into protected network segments.

**Example policy:**

```
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: restricted-network-vm-policy-1
spec:
  module: harvester-restricted-network-vm:20
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

**Specifications:**

1. You should be able to create a VM with any of the specific combinations there
2. You should not be able to create a VM from any namespace or network that is in that list, but the exact combination is not in the list.
3. Any namespace or network that is not on the list is not restricted

**Examples:**

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
