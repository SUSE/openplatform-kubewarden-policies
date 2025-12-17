#!/usr/bin/env bats

@test "accept because 12gb mig requested and a 12gb mig is in the ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-mig-12gb.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-mig-12gb.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "reject because 12gb mig requested and a 24gb mig is in the ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-mig-12gb.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-mig-24gb.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*MIG Partition 'nvidia.com/mig-1g.12gb' is not allowed for namespace: 'default'.*")" -ne 0 ]
}

@test "accept because 24gb mig requested and a 24gb mig is in the ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-mig-24gb.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-mig-24gb.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "reject because 24gb mig requested and a 12gb mig is in the ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-mig-24gb.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-mig-12gb.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*MIG Partition 'nvidia.com/mig-2g.24gb' is not allowed for namespace: 'default'.*")" -ne 0 ]
}

@test "reject because 12gb mig requested and there is no ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-mig-12gb.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-no-mig.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*MIG Partition 'nvidia.com/mig-1g.12gb' is not allowed for namespace: 'default'.*")" -ne 0 ]
}

@test "reject because 24gb mig requested and there is no ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-mig-24gb.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-no-mig.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
  [ "$(expr "$output" : ".*MIG Partition 'nvidia.com/mig-2g.24gb' is not allowed for namespace: 'default'.*")" -ne 0 ]
}

@test "accept because no mig requested and a 12gb mig is in the ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-no-mig.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-mig-12gb.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "accept because no mig requested and a 24gb mig is in the ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-no-mig.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-mig-24gb.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "accept because no mig requested and a 24gb mig is no ResourceQuota" {
  run kwctl run annotated-policy.wasm --request-path test_data/pod-no-mig.json --allow-context-aware --replay-host-capabilities-interactions test_data/session-no-mig.yaml
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}
