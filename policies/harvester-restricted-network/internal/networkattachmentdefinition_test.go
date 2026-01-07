package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	jsonString := `
{
    "apiVersion": "k8s.cni.cncf.io/v1",
    "kind": "NetworkAttachmentDefinition",
    "metadata": {
        "name": "network-1",
        "namespace": "test-restricted-1-network-1"
    },
    "spec": {
        "config": "{\"cniVersion\":\"0.3.1\",\"name\":\"network-1\",\"type\":\"bridge\",\"bridge\":\"mgmt-br\",\"promiscMode\":true,\"vlan\":1337,\"ipam\":{}}"
    }
}
`
	var networkDefinition NetworkAttachmentDefinition
	err := json.Unmarshal([]byte(jsonString), &networkDefinition)
	assert.NoError(t, err)
}
