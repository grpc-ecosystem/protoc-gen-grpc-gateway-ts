#!/bin/bash

function runTest {
  ORIG_NAME=${1:="false"}
  CONFIG_NAME=${2:="karma.conf.js"}
  ./scripts/gen-protos.sh $ORIG_NAME
  go run ./ -orig=$ORIG_NAME &
  pid=$!

  USE_PROTO_NAMES=$ORIG_NAME ./node_modules/.bin/karma start $CONFIG_NAME
  pkill -P $pid
}


