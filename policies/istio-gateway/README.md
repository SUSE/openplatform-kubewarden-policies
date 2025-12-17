# Kubewarden policy

## Description

This policy will restrict the usage of the `Gateway` object configured on top of a dedicated instance of istio proxy only to `VirtualService` object coming from a set of namespaces.

**Example policy:**

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: istio-gw-policy-1
spec:
  module: harbor.op-prg2-0-dev-ingress.op.suse.org/policy-istio-gateway/policy-istio-gateway:0.1.0
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

**Specifications:**

1. You should be able to create a Gateway only on specific namespaces for specific hosts and destination_hosts if defined, otherwise the `*` wildcard will allow `all`.
2. You should not be able to create a Gateway without specifying a valid namespace.

**Examples:**

```json
{
  "gatewayRestrictions": {
    "gateway-01": {
      "ns-01": [
        {
          "hostnames": ["*"],
          "destination_hosts": ["*"],
          "port": 443,
          "protocol": "https"
        }
      ]
    },
    "gateway-02": {
      "ns-02": [
        {
          "hostnames": ["hostname a"],
          "destination_hosts": ["servicename a", "servicename b"],
          "port": 80,
          "protocol": "http"
        },
        {
          "hostnames": ["hostname b"],
          "destination_hosts": ["servicename a", "servicename c"],
          "port": 443,
          "protocol": "https"
          }
      ]
    },
    "gateway-03": {
      "ns-03": [
        {
          "hostnames": ["hostname a"],
          "destination_hosts": ["*"],
          "port": 443,
          "protocol": "https"
          }
      ]
    }
  }
}
```
