[![Kubewarden Policy Repository](https://github.com/kubewarden/community/blob/main/badges/kubewarden-policies.svg)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#policy-scope)
[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# pod-mig-partitions

> [!NOTE]
> This project is meant to work with [NVIDIA GPU Operator](https://github.com/NVIDIA/gpu-operator).

With the NVIDIA GPU Operator, pods request MIG partitions with resource requests.
This policy ensures that a pod can only request a MIG partition within a namespace's [ResourceQuota](https://kubernetes.io/docs/concepts/policy/resource-quotas/).

## Example

The policy doesn't require any configuration, so you just need to add it to a Kubewarden policy server.

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: pod-mig-partitions
spec:
  module: registry://ghcr.io/suse/openplatform-kubewarden-policies/pod-mig-partitions:latest
  rules:
    - apiGroups: [""]
      apiVersions: ["v1"]
      resources: ["pods"]
      operations: ["CREATE", "UPDATE"]
  settings:
    mutating: false
    policyServer: default
```

With the policy active, if a pod tried to create or update a pod, adding a MIG partition, this policy should deny the change.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ollama
spec:
  containers:
    - image: dp.apps.rancher.io/containers/ollama:latest
      name: ollama
      resources:
        limits:
          cpu: '8'
          memory: 32Gi
          nvidia.com/mig-1g.12gb: '1'
        requests:
          cpu: '8'
          memory: 32Gi
          nvidia.com/mig-1g.12gb: '1'
```

To get the pod to deploy, would need to add a ResourceQuota with the requested resource.

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: gpu-quota
spec:
  hard:
    requests.nvidia.com/mig-1g.12gb: '1
```

Now, if the above pod requests the same MIG partition, it should be allowed. The pod-mig-partitions policy will see that `nvidia.com/mig-1g.12gb` is in the namespace's ResourceQuota and allow the change.
If the pod instead requests `requests.nvidia.com/mig-2g.24gb`, the policy would deny the change because that MIG partition is not in the ResourceQuota.
However, the policy doesn't concern itself with how many MIG partitions are in the request, instead,
Kubernetes ensures that the Pod doesn't exceed the namespace's ResourceQuotas.
