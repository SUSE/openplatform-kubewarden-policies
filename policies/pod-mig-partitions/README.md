# policy-pod-mig-partitions

This policy guards against Pods attaching MIG Partitions without a ResourceQuota.

**Example policy:**

```
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: pod-mig-partitions
spec:
  module: harbor.op-prg2-0-dev-ingress.op.suse.org/policy-pod-mig-partitions/policy-pod-mig-partitions:0.1.0
  rules:
    - apiGroups: [""]
      apiVersions: ["v1"]
      resources: ["pods"]
      operations: ["CREATE", "UPDATE"]
  settings:
    mutating: false  # or true if your policy mutates resources
    policyServer: default
```
