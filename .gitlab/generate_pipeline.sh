#!/usr/bin/env bash
set -e

OUTPUT_FILE=".gitlab/generated-pipeline.yml"

cat <<EOF > "$OUTPUT_FILE"
stages:
  - pipeline

.pipeline_template:
  stage: pipeline
  variables:
    GITLAB_TOKEN: \$CI_JOB_TOKEN
    GITLAB_USERNAME: gitlab-ci-token
    PARENT_PIPELINE_ID: "$CI_PIPELINE_ID"
    POLICY_WORKING_DIR: "${POLICY_WORKING_DIR}"
    POLICY_LANGUAGE: "${POLICY_LANGUAGE}"
    POLICY_RUST_PACKAGE: "${POLICY_RUST_PACKAGE}"
    POLICY_ID: "${POLICY_ID}"
    POLICY_VERSION: "${POLICY_VERSION}"
    POLICY_BASENAME: "${POLICY_BASENAME}"
EOF

if [ "${POLICY_LANGUAGE}" = "rust" ]; then

  cat <<EOF >> "$OUTPUT_FILE"
rust_pipeline:
  extends: .pipeline_template
  trigger:
    include: .gitlab/rust.gitlab-ci.yml
    strategy: depend
EOF

elif [ "${POLICY_LANGUAGE}" = "go" ]; then

 cat <<EOF >> "$OUTPUT_FILE"
go_pipeline:
  extends: .pipeline_template
  trigger:
    include: .gitlab/go.gitlab-ci.yml
    strategy: depend
EOF

fi
