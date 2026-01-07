package main

type vmMetadata struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type multus struct {
	NetworkName string `json:"networkName"`
}

type vmNetwork struct {
	Multus multus `json:"multus"`
}

type vmNetworks struct {
	Networks []vmNetwork `json:"networks"`
}

type vmTemplate struct {
	Spec vmNetworks `json:"spec"`
}

type vmPayloadSpec struct {
	Template vmTemplate `json:"template"`
}

type virtualMachine struct {
	Metadata vmMetadata    `json:"metadata"`
	Spec     vmPayloadSpec `json:"spec"`
}
