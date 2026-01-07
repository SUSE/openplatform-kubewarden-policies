package domain

type Metadata struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type PCIDevice struct {
	DeviceName string `json:"deviceName"`
	Name       string `json:"name"`
}

type Devices struct {
	GPUS        []PCIDevice `json:"gpus"`
	HostDevices []PCIDevice `json:"hostDevices"`
}

type Domain struct {
	Devices Devices `json:"devices"`
}

type VirtualMachineTemplateSpec struct {
	Domain Domain `json:"domain"`
}

type VirtualMachineSpecTemplate struct {
	Spec VirtualMachineTemplateSpec `json:"spec"`
}

type VirtualMachineSpec struct {
	Template VirtualMachineSpecTemplate `json:"template"`
}

type VirtualMachine struct {
	Metadata Metadata           `json:"metadata"`
	Spec     VirtualMachineSpec `json:"spec"`
}
