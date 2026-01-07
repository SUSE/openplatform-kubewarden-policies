package internal

import "encoding/json"

type Config struct {
	CniVersion  string      `json:"cniVersion"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Bridge      string      `json:"bridge"`
	PromiscMode bool        `json:"promiscMode"`
	VLAN        int         `json:"vlan"`
	IPAM        interface{} `json:"ipam"`
}

type Metadata struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type Spec struct {
	Config Config `json:"config"`
}

// UnmarshalJSON overrides the default implementation of JSON unmarshalling.
// This is required because the config object is a JSON string and not a JSON object.
func (c *Config) UnmarshalJSON(data []byte) error {
	// unmarshal the data to a string
	var jsonString string
	_ = json.Unmarshal(data, &jsonString)

	// unmarshal the json string to an object
	// Create a new type to avoid infinite recursion
	type C Config
	return json.Unmarshal([]byte(jsonString), (*C)(c))
}

type NetworkAttachmentDefinition struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
}
