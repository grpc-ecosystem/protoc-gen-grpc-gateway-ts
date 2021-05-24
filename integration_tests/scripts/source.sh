#!/bin/bash

function runTest {
  ORIG_NAME=${1:="false"}
  CONFIG_NAME=${2:="karma.conf.js"}
  ./scripts/gen-protos.sh $ORIG_NAME true
  go run ./ -orig=$ORIG_NAME &
  pid=$!

  USE_PROTO_NAMES=$ORIG_NAME ./node_modules/.bin/karma start $CONFIG_NAME
  TEST_EXIT=$?
  if [[ $TEST_EXIT -ne 0 ]]; then
    pkill -P $pid
    exit $TEST_EXIT
  fi

  pkill -P $pid
}


