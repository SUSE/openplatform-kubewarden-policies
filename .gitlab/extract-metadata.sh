#!/usr/bin/env bash
set -e

policy_name=$( echo ${CI_COMMIT_TAG} | sed 's/\(.*\)\/\(.*\)$/\1/' )
policy_working_dir=$( find policies -type d -name "$policy_name" )
echo "POLICY_WORKING_DIR=$policy_working_dir"

if [ ! -d "$policy_working_dir" ]; then
  echo "$policy_working_dir does not exist, policy not found";
  exit 1;
fi

policy_ociUrl=$(yq -r '.annotations."io.kubewarden.policy.ociUrl"' "${policy_working_dir}/metadata.yml")
policy_version=$(yq -r '.annotations."io.kubewarden.policy.version"' "${policy_working_dir}/metadata.yml")
policy_id=${policy_ociUrl##*/}
policy_basename=$(basename ${policy_working_dir})
policy_language=""
policy_rust_package=""

if [ -f "${policy_working_dir}/Cargo.toml" ]; then
  policy_language="rust"
  policy_rust_package=$(sed  -n 's,^name = \"\(.*\)\",\1,p' "${policy_working_dir}/Cargo.toml")
  if [ "$policy_rust_package" == "" ]; then
    echo "cannot get rust policy ${POLICY_WORKING_DIR} package name";
    exit 1;
  fi
else
  # Currently this repository supports go and rust policies only
  policy_language="go"
fi

echo "POLICY_LANGUAGE=$policy_language"
echo "POLICY_RUST_PACKAGE=$policy_rust_package"
echo "POLICY_ID=$policy_id"
echo "POLICY_VERSION=$policy_version"
echo "POLICY_BASENAME=$policy_basename"

echo "POLICY_WORKING_DIR=$policy_working_dir" >> build.env
echo "POLICY_LANGUAGE=$policy_language" >> build.env
echo "POLICY_RUST_PACKAGE=$policy_rust_package" >> build.env
echo "POLICY_ID=$policy_id" >> build.env
echo "POLICY_VERSION=$policy_version" >> build.env
echo "POLICY_BASENAME=$policy_basename" >> build.env
