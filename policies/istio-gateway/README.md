[![Kubewarden Policy Repository](https://github.com/kubewarden/community/blob/main/badges/kubewarden-policies.svg)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#policy-scope)
[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# istio-gateway

> [!NOTE]
> This policy is meant to work with link:https://istio.io/[Istio], but not does not protect resources from its [Gateway API](https://istio.io/latest/docs/tasks/traffic-management/ingress/gateway-api/) implementation.

This policy protects shared Istio Gateway resources by watching changes to VirtualService resources.
For configured Gateway resources, it will ensure that VirtualService resources are correctly configured.

## Settings

| Field                                                                           | Description                                  |
|---------------------------------------------------------------------------------|----------------------------------------------|
| gatewayRestirctions <br/> map[string, [gatewayRestriction](#gatewayRestriction) | A list of Istio Gateway objects to restrict. |

### GatewayRestriction

| Field                                                | Description                 |
|------------------------------------------------------|-----------------------------|
| namespaces <br> map[string, [namespace](#namespace)] | A map of namespace objects. |

### Namespace

| Field                            | Description                                                                      |
|----------------------------------|----------------------------------------------------------------------------------|
| hostnames <br/> string[]         | A list of hostnames for the VirtualService.                                      |
| port <br/> int                   | The port for the VirtualService. The default value 0 means any.                  |
| protocol <br/> string            | The protocol for the VirtualService. The default value (empty string) means any. |
| destination_hosts <br/> string[] | The destination hosts for the VirtualService.                                    |

## Specifications

1. You should be able to create a Gateway only on specific namespaces for specific hosts and destination_hosts if defined, otherwise the `*` wildcard will allow `all`.
2. You should not be able to create a Gateway without specifying a valid namespace.

## Example

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: istio-gw-policy-1
spec:
  module: registry://ghcr.io/suse/openplatform-kubewarden-policies/istio-gateway:latest
  rules:
    - apiGroups: ["networking.istio.io"]
      apiVersions: ["v1"]
      resources: ["virtualservices"]
      operations: ["CREATE", "UPDATE"]
  settings:
    gatewayRestrictions:
      "gateway01":
        "ns-1":
          - hostnames: []
            destination_hosts: []
      "gateway02":
        "ns-2":
          - hostnames: ["hostname a"]
            port: "80"
            protocol: "http"
            destination_hosts: ["servicename a", "servicename b"]
          - hostnames: ["hostname b"]
            port: "443"
            protocol: "https"
            destination_hosts: ["servicename a", "servicename c"]
        "ns-3":
          - hostnames: ["hostname c"]
            port: "443"
            protocol: "https"
            destination_hosts: []
  mutating: false
  policyServer: default
```
