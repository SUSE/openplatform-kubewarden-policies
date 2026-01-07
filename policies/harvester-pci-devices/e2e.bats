#!/usr/bin/env bats

@test "accept because the gpu is bound to the namespace" {
  run kwctl run annotated-policy.wasm -r test_data/virtualmachine-gpu.json --settings-json '{"namespaceDeviceBindings": [{"namespace": "default", "device": "tekton27a-000001010"}]}'
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "accept because the pci device is bound to the namespace" {
  run kwctl run annotated-policy.wasm -r test_data/virtualmachine-pci.json --settings-json '{"namespaceDeviceBindings": [{"namespace": "default", "device": "tekton27a-000001010"}]}'
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "accept because a gpu is not requested" {
  run kwctl run annotated-policy.wasm -r test_data/virtualmachine-no-gpu.json
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "reject because gpu is not bound to namespace" {
  run kwctl run annotated-policy.wasm -r test_data/virtualmachine-gpu.json --settings-json '{"namespaceDeviceBindings": [{"namespace": "foo", "device": "test-gpu"}]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*PCI DEVICE 'tekton27a-000001010' is not allowed for namespace: 'default'.*")" -ne 0 ]
}

@test "reject because pci device is not bound to namespace" {
  run kwctl run annotated-policy.wasm -r test_data/virtualmachine-pci.json --settings-json '{"namespaceDeviceBindings": [{"namespace": "foo", "device": "test-gpu"}]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*PCI DEVICE 'tekton27a-000001010' is not allowed for namespace: 'default'.*")" -ne 0 ]
}

@test "reject because gpu is not bound to another namespace" {
  run kwctl run annotated-policy.wasm -r test_data/virtualmachine-gpu.json --settings-json '{"namespaceDeviceBindings": [{"namespace": "foobar", "device": "tekton27a-000001010"}]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*PCI DEVICE 'tekton27a-000001010' is not allowed for namespace: 'default'.*")" -ne 0 ]
}

@test "reject because pci device is not bound to another namespace" {
  run kwctl run annotated-policy.wasm -r test_data/virtualmachine-pci.json --settings-json '{"namespaceDeviceBindings": [{"namespace": "foobar", "device": "tekton27a-000001010"}]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*PCI DEVICE 'tekton27a-000001010' is not allowed for namespace: 'default'.*")" -ne 0 ]
}
