#!/usr/bin/env bats

POLICY_WASM_PATH="annotated-policy.wasm"
SETTINGS_FILE="test_data/settings.json"

@test "accepts a VS from an allowed namespace using a restricted gateway" {
  run kwctl run \
    --settings-path ${SETTINGS_FILE} \
    --request-path "test_data/vs-01.json" \
    ${POLICY_WASM_PATH}

  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*true')" -ne 0 ]
}

@test "rejects a VS from a not allowed namespace using a restricted gateway" {
  run kwctl run \
    --settings-path ${SETTINGS_FILE} \
    --request-path "test_data/vs-02.json" \
    ${POLICY_WASM_PATH}

  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
}

@test "rejects a VS from an allowed namespace using a restricted gateway but a not allowed service" {
  run kwctl run \
    --settings-path "test_data/settings2.json" \
    --request-path "test_data/vs-03.json" \
    ${POLICY_WASM_PATH}

  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ "$(expr "$output" : '.*allowed.*false')" -ne 0 ]
}
