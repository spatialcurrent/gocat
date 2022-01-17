#!/bin/bash

# =================================================================
#
# Copyright (C) 2022 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export testdata_local="${DIR}/../testdata"

export "testdata_http"="${GOCAT_TESTDATA_HTTP:-}"

export "testdata_s3"="${GOCAT_TESTDATA_S3:-}"

testHelp() {
  "${DIR}/../bin/gocat" --help
}

#
# Local
#

testLocalStdin() {
  local input="hello world"
  local expected='hello world'
  local output=$(echo "${input}" | "${DIR}/../bin/gocat" -)
  assertEquals "unexpected output" "${expected}" "${output}"
}

testLocalFile() {
  local input="hello world"
  local expected='hello world'
  local output=$("${DIR}/../bin/gocat" "${testdata_local}/doc.txt")
  assertEquals "unexpected output" "${expected}" "${output}"
}

testLocalFiles() {
  local input="hello world"
  local expected='hello world\nhello world'
  local output=$("${DIR}/../bin/gocat" "${testdata_local}/doc.txt" "${testdata_local}/doc.txt")
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

#
# HTTP
#

testHttpFile() {
  if [[ ! -z "${testdata_http}" ]]; then
    local input="hello world"
    local expected='hello world'
    local output=$("${DIR}/../bin/gocat" "${testdata_http}/doc.txt")
    assertEquals "unexpected output" "${expected}" "${output}"
  else
    echo "* skipping"
  fi
}

testHttpFiles() {
  if [[ ! -z "${testdata_http}" ]]; then
    local input="hello world\nhello world"
    local expected='hello world'
    local output=$("${DIR}/../bin/gocat" "${testdata_http}/doc.txt" "${testdata_http}/doc.txt")
    assertEquals "unexpected output" "${expected}" "${output}"
  else
    echo "* skipping"
  fi
}

#
# S3
#

testS3File() {
  if [[ ! -z "${testdata_s3}" ]]; then
    local input="hello world"
    local expected='hello world'
    local output=$("${DIR}/../bin/gocat" "${testdata_s3}/doc.txt")
    assertEquals "unexpected output" "${expected}" "${output}"
  else
    echo "* skipping"
  fi
}

testS3Files() {
  if [[ ! -z "${testdata_s3}" ]]; then
    local input="hello world\nhello world"
    local expected='hello world'
    local output=$("${DIR}/../bin/gocat" "${testdata_s3}/doc.txt" "${testdata_s3}/doc.txt")
    assertEquals "unexpected output" "${expected}" "${output}"
  else
    echo "* skipping"
  fi
}


oneTimeSetUp() {
  echo "Using temporary directory at ${SHUNIT_TMPDIR}"
  echo "Reading testdata from ${testdata_local}"
}

oneTimeTearDown() {
  echo "Tearing Down"
}

# Load shUnit2.
. "${DIR}/shunit2"
