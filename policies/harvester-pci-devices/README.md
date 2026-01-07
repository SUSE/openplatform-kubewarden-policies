# policy-pci-devices

This policy guards against VMs attaching PCI Devices (e.g., GPUs) without permission.

**Example policy:**

```
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

**Specifications:**

1. You should be able to create a VM without a PCI Device
2. You should not be able to bind a VM to a PCI Device not allocated to its namespace.

**Examples:**

| Namespace        | Device              | Result |
|------------------|---------------------|--------|
| test-ns-1        | tekton27a-000001010 | ALLOW  |
| test-ns-2        | tekton27b-000001010 | ALLOW  |
| test-ns-1        | tekton27b-000001010 | REJECT |
| test-ns-2        | tekton27a-000001010 | REJECT |
| random-namespace | tekton27a-000001010 | REJECT |
| random-namespace | tekton27b-000001010 | REJECT |
